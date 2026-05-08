package monitoring

import (
	"testing"
)

func TestMonitoringDashboardDiffSuppress_emptyContainersDroppedByApi(t *testing.T) {
	// The Monitoring Dashboard API drops empty maps and empty arrays on
	// round-trip. A user config that includes "labels": {} or
	// "dashboardFilters": [] must not show a permanent diff against the
	// state read back from the API (which omits those keys).
	cases := map[string]struct {
		old, new string
		want     bool
	}{
		"empty_labels_in_new_only": {
			old:  `{"displayName":"d"}`,
			new:  `{"displayName":"d","labels":{}}`,
			want: true,
		},
		"empty_array_in_new_only": {
			old:  `{"displayName":"d"}`,
			new:  `{"displayName":"d","dashboardFilters":[]}`,
			want: true,
		},
		"both_empty_containers_in_new": {
			old:  `{"displayName":"d","mosaicLayout":{"columns":12}}`,
			new:  `{"displayName":"d","labels":{},"dashboardFilters":[],"mosaicLayout":{"columns":12}}`,
			want: true,
		},
		"nested_empty_array_in_widget": {
			old:  `{"mosaicLayout":{"tiles":[{"widget":{"xyChart":{"dataSets":[{"plotType":"LINE"}]}}}]}}`,
			new:  `{"mosaicLayout":{"tiles":[{"widget":{"xyChart":{"chartOptions":{},"dataSets":[{"breakdowns":[],"dimensions":[],"measures":[],"plotType":"LINE"}]}}}]}}`,
			want: true,
		},
		"non_empty_label_value_must_diff": {
			old:  `{"displayName":"d"}`,
			new:  `{"displayName":"d","labels":{"env":"prod"}}`,
			want: false,
		},
		"display_name_change_must_diff": {
			old:  `{"displayName":"old"}`,
			new:  `{"displayName":"new"}`,
			want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := monitoringDashboardDiffSuppress("dashboard_json", tc.old, tc.new, nil)
			if got != tc.want {
				t.Errorf("monitoringDashboardDiffSuppress(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.want)
			}
		})
	}
}

func TestMonitoringDashboardDiffSuppress_typeMismatchSafe(t *testing.T) {
	// Regression: removeComputedKeys used to panic on type mismatch
	// (e.g. old has key X as a slice, new has X as a map). After the
	// type-assertion guards, mismatches must just keep the diff visible
	// instead of panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("monitoringDashboardDiffSuppress panicked: %v", r)
		}
	}()
	old := `{"x":[1,2,3]}`
	new := `{"x":{"a":1}}`
	got := monitoringDashboardDiffSuppress("dashboard_json", old, new, nil)
	if got {
		t.Errorf("expected diff (got suppress=true) for type-mismatched values")
	}
}
