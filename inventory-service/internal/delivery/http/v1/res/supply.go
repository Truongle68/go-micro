package res

import (
	"inventory-service/internal/domain"
)

type ProductAssignedRes struct {
	Summary Summary                         `json:"summary"`
	Results []domain.SupplyAssignmentResult `json:"results"`
}

type Summary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

func ToProductAssignedRes(results []domain.SupplyAssignmentResult) ProductAssignedRes {
	succeeded, failed := countResult(results)

	summary := Summary{
		Total:     len(results),
		Succeeded: succeeded,
		Failed:    failed,
	}

	for _, r := range results {
		r.MarshalJSON()
	}

	return ProductAssignedRes{
		Summary: summary,
		Results: results,
	}
}

func countResult(results []domain.SupplyAssignmentResult) (succeeded, failed int) {
	for _, r := range results {
		if r.Success {
			succeeded++
		} else {
			failed++
		}
	}
	return succeeded, failed
}
