package tags

import "context"

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	DummyTags    = &dummyTags{}
)

type Tags interface {
	Has(key string) bool
	Values() map[string]interface{}
	Set(key string, value interface{}) Tags
}

type mapTags struct {
	values map[string]interface{}
}

func (t *mapTags) Has(key string) bool {
	_, ok := t.values[key]
	return ok
}

func (t *mapTags) Values() map[string]interface{} {
	return t.values
}

func (t *mapTags) Set(key string, value interface{}) Tags {
	t.values[key] = value
	return t
}

type dummyTags struct{}

func (t *dummyTags) Has(key string) bool {
	return false
}

func (t *dummyTags) Values() map[string]interface{} {
	return nil
}

func (t *dummyTags) Set(key string, value interface{}) Tags {
	return t
}

func NewTags() Tags {
	return &mapTags{values: make(map[string]interface{})}
}

func Extract(ctx context.Context) Tags {
	tags, ok := ctx.Value(ctxMarkerKey).(Tags)
	if !ok {
		return DummyTags
	}
	return tags
}

func SetInContext(ctx context.Context, tags Tags) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, tags)
}
