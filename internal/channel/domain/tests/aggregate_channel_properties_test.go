package channel_tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	channel_domain "vault-app/internal/channel/domain"
)

func testChannelPropertyOne() channel_domain.ChannelProperty {
	return channel_domain.ChannelProperty{
		Key:   "review_required",
		Value: "true",
	}
}

func testChannelPropertyTwo() channel_domain.ChannelProperty {
	return channel_domain.ChannelProperty{
		Key:   "download_allowed",
		Value: "false",
	}
}

func testChannelPropertyThree() channel_domain.ChannelProperty {
	return channel_domain.ChannelProperty{
		Key:   "required_role",
		Value: "lead-engineer",
	}
}

func TestChannel_AddChannelProperty_GetChannelProperty(t *testing.T) {
	channel := newTestChannel()
	property := testChannelPropertyOne()

	result := channel.AddChannelProperty(property)

	require.Len(t, result.Properties, 1)
	require.Equal(t, property, result.Properties[0])

	found, ok := result.GetChannelProperty("review_required")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, property, *found)
}

func TestChannel_GetChannelProperty_NotFound(t *testing.T) {
	channel := newTestChannel()

	found, ok := channel.GetChannelProperty("unknown-key")

	require.False(t, ok)
	require.Nil(t, found)
}

func TestChannel_UpdateChannelProperty(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddChannelProperty(testChannelPropertyOne())

	updated := channel_domain.ChannelProperty{
		Key:   "review_required",
		Value: "false",
	}

	updatedResult := channel.UpdateChannelProperty(updated)

	require.True(t, updatedResult)

	found, ok := channel.GetChannelProperty("review_required")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, "review_required", found.Key)
	require.Equal(t, "false", found.Value)
}

func TestChannel_UpdateChannelProperty_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddChannelProperty(testChannelPropertyOne())

	updated := channel_domain.ChannelProperty{
		Key:   "unknown-key",
		Value: "unknown-value",
	}

	updatedResult := channel.UpdateChannelProperty(updated)

	require.False(t, updatedResult)
	require.Len(t, channel.Properties, 1)
	require.Equal(t, testChannelPropertyOne(), channel.Properties[0])
}

func TestChannel_RemoveChannelProperty(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddChannelProperty(testChannelPropertyOne())
	channel = channel.AddChannelProperty(testChannelPropertyTwo())

	removed := channel.RemoveChannelProperty("review_required")

	require.True(t, removed)
	require.Len(t, channel.Properties, 1)
	require.Equal(t, "download_allowed", channel.Properties[0].Key)

	_, found := channel.GetChannelProperty("review_required")
	require.False(t, found)
}

func TestChannel_RemoveChannelProperty_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddChannelProperty(testChannelPropertyOne())

	removed := channel.RemoveChannelProperty("unknown-key")

	require.False(t, removed)
	require.Len(t, channel.Properties, 1)
	require.Equal(t, testChannelPropertyOne(), channel.Properties[0])
}

func TestChannel_ListChannelProperties_ReturnsCopy(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddChannelProperty(testChannelPropertyOne())

	properties := channel.ListChannelProperties()

	require.Len(t, properties, 1)

	properties[0].Value = "changed-outside-channel"
	properties = append(properties, testChannelPropertyTwo())

	require.Len(t, channel.Properties, 1)

	original, ok := channel.GetChannelProperty("review_required")

	require.True(t, ok)
	require.Equal(t, "true", original.Value)
}

func TestChannel_ChannelPropertyCRUD_MultipleProperties(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddChannelProperty(testChannelPropertyOne())
	channel = channel.AddChannelProperty(testChannelPropertyTwo())
	channel = channel.AddChannelProperty(testChannelPropertyThree())

	require.Len(t, channel.Properties, 3)

	review, found := channel.GetChannelProperty("review_required")
	require.True(t, found)
	require.Equal(t, "true", review.Value)

	download, found := channel.GetChannelProperty("download_allowed")
	require.True(t, found)
	require.Equal(t, "false", download.Value)

	role, found := channel.GetChannelProperty("required_role")
	require.True(t, found)
	require.Equal(t, "lead-engineer", role.Value)
}

func TestChannel_ChannelPropertyCRUD_EmptyChannel(t *testing.T) {
	channel := newTestChannel()

	require.Empty(t, channel.ListChannelProperties())

	property, found := channel.GetChannelProperty("review_required")

	require.False(t, found)
	require.Nil(t, property)

	require.False(t, channel.UpdateChannelProperty(
		testChannelPropertyOne(),
	))

	require.False(t, channel.RemoveChannelProperty(
		"review_required",
	))
}

func TestChannel_AddChannelPropertyIfAbsent(t *testing.T) {
	channel := newTestChannel()

	firstAdded := channel.AddChannelPropertyIfAbsent(
		testChannelPropertyOne(),
	)

	secondAdded := channel.AddChannelPropertyIfAbsent(
		testChannelPropertyOne(),
	)

	require.True(t, firstAdded)
	require.False(t, secondAdded)
	require.Len(t, channel.Properties, 1)
}
