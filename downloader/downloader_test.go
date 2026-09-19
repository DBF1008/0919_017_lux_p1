package downloader

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/iawia002/lux/extractors"
)

func writeTestPartFile(t *testing.T, filePath string, part *FilePartMeta, data []byte) {
	t.Helper()
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close() // nolint
	if err := binary.Write(file, binary.LittleEndian, part); err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}
}

func TestMergeMultiPart(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.mp4")
	parts := []*FilePartMeta{
		{Index: 0, Start: 0, End: 4, Cur: 5},
		{Index: 1, Start: 5, End: 9, Cur: 10},
	}
	writeTestPartFile(t, filePartPath(filePath, parts[0]), parts[0], []byte("hello"))
	writeTestPartFile(t, filePartPath(filePath, parts[1]), parts[1], []byte("world"))

	if err := mergeMultiPart(filePath, parts); err != nil {
		t.Fatal(err)
	}
	merged, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(merged, []byte("helloworld")) {
		t.Fatalf("unexpected merged content: %q", merged)
	}
	// part files should be cleaned up after merging
	for _, part := range parts {
		if _, err := os.Stat(filePartPath(filePath, part)); !os.IsNotExist(err) {
			t.Fatalf("part file %s should have been removed", filePartPath(filePath, part))
		}
	}
}

func TestMergeMultiPartMissingPartFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.mp4")
	parts := []*FilePartMeta{
		{Index: 0, Start: 0, End: 4, Cur: 5},
		{Index: 1, Start: 5, End: 9, Cur: 10},
	}
	writeTestPartFile(t, filePartPath(filePath, parts[0]), parts[0], []byte("hello"))
	// parts[1] has no corresponding file, mergeMultiPart should return an
	// error instead of panicking.
	if err := mergeMultiPart(filePath, parts); err == nil {
		t.Fatal("expected an error for missing part file, got nil")
	}
}

func TestDownload(t *testing.T) {
	testCases := []struct {
		name string
		data *extractors.Data
	}{
		{
			name: "normal test",
			data: &extractors.Data{
				Site:  "douyin",
				Title: "test",
				Type:  extractors.DataTypeVideo,
				URL:   "https://www.douyin.com",
				Streams: map[string]*extractors.Stream{
					"default": {
						ID: "default",
						Parts: []*extractors.Part{
							{
								URL:  "https://aweme.snssdk.com/aweme/v1/playwm/?video_id=v0200f9a0000bc117isuatl67cees890&line=0",
								Size: 4927877,
								Ext:  "mp4",
							},
						},
					},
				},
			},
		},
		{
			name: "multi-stream test",
			data: &extractors.Data{
				Site:  "douyin",
				Title: "test2",
				Type:  extractors.DataTypeVideo,
				URL:   "https://www.douyin.com",
				Streams: map[string]*extractors.Stream{
					"miaopai": {
						ID: "miaopai",
						Parts: []*extractors.Part{
							{
								URL:  "https://txycdn.miaopai.com/stream/KwR26jUGh2ySnVjYbQiFmomNjP14LtMU3vi6sQ__.mp4?ssig=6594aa01a78e78f50c65c164d186ba9e&time_stamp=1537070910786",
								Size: 4011590,
								Ext:  "mp4",
							},
						},
						Size: 4011590,
					},
					"douyin": {
						ID: "douyin",
						Parts: []*extractors.Part{
							{
								URL:  "https://aweme.snssdk.com/aweme/v1/playwm/?video_id=v0200f9a0000bc117isuatl67cees890&line=0",
								Size: 4927877,
								Ext:  "mp4",
							},
						},
						Size: 4927877,
					},
				},
			},
		},
		{
			name: "image test",
			data: &extractors.Data{
				Site:  "bcy",
				Title: "bcy image test",
				Type:  extractors.DataTypeImage,
				URL:   "https://www.bcyimg.com",
				Streams: map[string]*extractors.Stream{
					"default": {
						ID: "default",
						Parts: []*extractors.Part{
							{
								URL:  "http://img5.bcyimg.com/coser/143767/post/c0j7x/0d713eb41a614053ac6a3b146914f6bc.jpg/w650",
								Size: 56107,
								Ext:  "jpg",
							},
							{
								URL:  "http://img9.bcyimg.com/coser/143767/post/c0j7x/d17e9b8587794d939a1363c5f715014b.jpg/w650",
								Size: 142100,
								Ext:  "jpg",
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		err := New(Options{}).Download(testCase.data)
		if err != nil {
			t.Error(err)
		}
	}
}
