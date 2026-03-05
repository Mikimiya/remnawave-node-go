package controller

import (
	"strings"

	appstats "github.com/xtls/xray-core/app/stats"
	"github.com/xtls/xray-core/features/stats"
)

func (c *StatsController) getStatsManager() stats.Manager {
	instance := c.core.Instance()
	if instance == nil {
		return nil
	}

	stmFeature := instance.GetFeature(stats.ManagerType())
	if stmFeature == nil {
		return nil
	}

	stm, ok := stmFeature.(stats.Manager)
	if !ok {
		return nil
	}

	return stm
}

func (c *StatsController) getConcreteStatsManager() *appstats.Manager {
	instance := c.core.Instance()
	if instance == nil {
		return nil
	}

	stmFeature := instance.GetFeature(stats.ManagerType())
	if stmFeature == nil {
		return nil
	}

	stm, ok := stmFeature.(*appstats.Manager)
	if !ok {
		return nil
	}

	return stm
}

func (c *StatsController) getCounterValue(stm stats.Manager, name string, reset bool) int64 {
	counter := stm.GetCounter(name)
	if counter == nil {
		return 0
	}
	value := counter.Value()
	if reset {
		counter.Set(0)
	}
	return value
}

func (c *StatsController) collectTrafficStats(stm *appstats.Manager, prefix string, reset bool) map[string]map[string]int64 {
	result := make(map[string]map[string]int64)

	stm.VisitCounters(func(name string, counter stats.Counter) bool {
		if !strings.HasPrefix(name, prefix) {
			return true
		}

		parts := strings.Split(name, ">>>")
		if len(parts) < 4 {
			return true
		}

		tag := parts[1]
		if parts[2] != "traffic" {
			return true
		}
		direction := parts[3]

		if result[tag] == nil {
			result[tag] = make(map[string]int64)
		}

		value := counter.Value()
		if reset {
			counter.Set(0)
		}

		result[tag][direction] = value
		return true
	})

	return result
}

func (c *StatsController) collectUserStats(stm *appstats.Manager, reset bool) map[string]*UserStats {
	userTraffic := make(map[string]*UserStats)

	stm.VisitCounters(func(name string, counter stats.Counter) bool {
		if !strings.HasPrefix(name, "user>>>") {
			return true
		}

		parts := strings.Split(name, ">>>")
		if len(parts) < 4 || parts[2] != "traffic" {
			return true
		}

		username := parts[1]
		direction := parts[3]

		value := counter.Value()
		if reset {
			counter.Set(0)
		}

		if userTraffic[username] == nil {
			userTraffic[username] = &UserStats{Username: username}
		}

		if direction == "uplink" {
			userTraffic[username].Uplink = value
		} else if direction == "downlink" {
			userTraffic[username].Downlink = value
		}

		return true
	})

	return userTraffic
}
