package permission

import "errors"

func CheckFeatureAccess(plan string, feature string) error {
	if features, ok := PlanFeatureMap[plan]; ok {
		if allowed, found := features[feature]; found {
			if allowed {
				return nil
			}
			return errors.New("feature not available on your plan")
		}
	}
	return errors.New("unknown plan or feature")
}

func GetFileSizeLimit(plan string) int64 {
	limits := map[string]int64{
		"free":    2 << 20,  // 2 MB
		"basic":   5 << 20,  // 5 MB
		"premium": 10 << 20, // 10 MB
		"custom":  20 << 20, // 20 MB
	}
	if limit, ok := limits[plan]; ok {
		return limit
	}
	return 2 << 20
}
