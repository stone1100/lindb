package operator

type groupingTagsCollect struct {
}

func (op *groupingTagsCollect) Execute() error {

	t.metadata.TagMetadata().CollectTagValues(t.tagKey.ID, t.tagValueIDs, t.tagValues)
}
