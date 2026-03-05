package xray

import (
	"fmt"
	"net"
	"strings"

	"github.com/xtls/xray-core/app/router"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/features/routing"
)

type routerWithRules interface {
	routing.Router
	AddRule(msg *serial.TypedMessage, shouldAppend bool) error
	RemoveRule(tag string) error
}

func (c *Core) getRouter() (routerWithRules, error) {
	c.mu.RLock()
	instance := c.instance
	c.mu.RUnlock()

	if instance == nil {
		return nil, fmt.Errorf("xray instance not running")
	}

	routerFeature := instance.GetFeature(routing.RouterType())
	if routerFeature == nil {
		return nil, fmt.Errorf("router feature not found")
	}

	r, ok := routerFeature.(routerWithRules)
	if !ok {
		return nil, fmt.Errorf("router does not support dynamic rule management")
	}

	return r, nil
}

func (c *Core) AddRoutingRule(ruleTag string, sourceIP string, outboundTag string) error {
	r, err := c.getRouter()
	if err != nil {
		return err
	}

	ip := net.ParseIP(sourceIP)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", sourceIP)
	}

	var ipBytes []byte
	var prefix uint32
	if ip4 := ip.To4(); ip4 != nil {
		ipBytes = ip4
		prefix = 32
	} else {
		ipBytes = ip.To16()
		prefix = 128
	}

	routerConfig := &router.Config{
		Rule: []*router.RoutingRule{
			{
				RuleTag: ruleTag,
				TargetTag: &router.RoutingRule_Tag{
					Tag: outboundTag,
				},
				SourceGeoip: []*router.GeoIP{
					{
						Cidr: []*router.CIDR{
							{
								Ip:     ipBytes,
								Prefix: prefix,
							},
						},
					},
				},
			},
		},
	}

	typedMsg := serial.ToTypedMessage(routerConfig)

	if err := r.AddRule(typedMsg, true); err != nil {
		return fmt.Errorf("failed to add routing rule: %w", err)
	}

	c.logger.WithField("ruleTag", ruleTag).WithField("sourceIP", sourceIP).
		WithField("outbound", outboundTag).Info("Added routing rule")

	return nil
}

func (c *Core) RemoveRoutingRule(ruleTag string) error {
	r, err := c.getRouter()
	if err != nil {
		return err
	}

	if err := r.RemoveRule(ruleTag); err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "empty tag") {
			c.logger.WithField("ruleTag", ruleTag).Warn("Rule not found, may already be removed")
			return nil
		}
		return fmt.Errorf("failed to remove routing rule: %w", err)
	}

	c.logger.WithField("ruleTag", ruleTag).Info("Removed routing rule")

	return nil
}
