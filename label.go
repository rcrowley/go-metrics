package metrics

type Label struct {
	Key   string
	Value string
}

func deepCopyLabels(labels []Label) []Label {
	copy := []Label{}
	for _, lbl := range labels {
		copy = append(copy, Label{
			Key:   lbl.Key,
			Value: lbl.Value,
		})
	}
	return copy
}
