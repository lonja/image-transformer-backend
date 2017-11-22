package api

import (
	"testing"

	"model"
	"io"
	"mime/multipart"
	"os"
	"bytes"
	"net/http"
	"time"
	"strconv"
)

const baseURL = "http://localhost:9898"

var client = &http.Client{}

func Test_handleRotation(t *testing.T) {
	go New().Start(9898)
	type args struct {
		files []string
		angle string
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantStatus int
	}{
		{
			name:       "empty",
			args:       args{},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty file",
			args: args{
				angle: "90",
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "1 file angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
				angle: "90",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "couple files angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
				},
				angle: "90",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			fw, err := w.CreateFormField("angle")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(tt.args.angle)); err != nil {
				t.Fail()
				return
			}
			for _, file := range tt.args.files {
				f, err := os.Open(file)
				if err != nil {
					t.Fail()
					return
				}
				defer f.Close()
				fw, err = w.CreateFormFile("file", file)
				if err != nil {
					t.Fail()
					return
				}
				if _, err = io.Copy(fw, f); err != nil {
					t.Fail()
					return
				}
			}
			w.Close()
			req, err := http.NewRequest("POST", baseURL+"/rotate", &b)
			if err != nil {
				t.Fail()
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			startTime := time.Now().UnixNano()
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano()
			if err != nil {
				t.Fail()
				return
			}
			t.Logf(`request "%v" processing duration: %v sec`, tt.name, float32(endTime-startTime)/1000000000)
			if resp.StatusCode != tt.wantStatus || err != nil {
				t.Errorf("rotation error, expected status %v, having status %v", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func Test_handleResize(t *testing.T) {
	go New().Start(9898)
	type args struct {
		files     []string
		width     string
		height    string
		keepRatio bool
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantStatus int
	}{
		{
			name:       "empty",
			args:       args{},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty file",
			args: args{
				height: "500px",
				width:  "500px",
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "1 file size",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
				height: "500px",
				width:  "500px",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "1 file size keep ratio",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
				height:    "500px",
				width:     "500px",
				keepRatio: true,
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "couple files angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
				},
				height: "500px",
				width:  "500px",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			fw, err := w.CreateFormField("width")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(tt.args.width)); err != nil {
				t.Fail()
				return
			}
			fw, err = w.CreateFormField("height")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(tt.args.height)); err != nil {
				t.Fail()
				return
			}
			fw, err = w.CreateFormField("keepRatio")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(strconv.FormatBool(tt.args.keepRatio))); err != nil {
				t.Fail()
				return
			}
			for _, file := range tt.args.files {
				f, err := os.Open(file)
				if err != nil {
					t.Fail()
					return
				}
				defer f.Close()
				fw, err = w.CreateFormFile("file", file)
				if err != nil {
					t.Fail()
					return
				}
				if _, err = io.Copy(fw, f); err != nil {
					t.Fail()
					return
				}
			}
			w.Close()
			req, err := http.NewRequest("POST", baseURL+"/resize", &b)
			if err != nil {
				t.Fail()
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			startTime := time.Now().UnixNano()
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano()
			t.Logf(`request "%v" processing duration: %v sec`, tt.name, float32(endTime-startTime)/1000000000)
			if err != nil {
				t.Fail()
				return
			}
			if resp.StatusCode != tt.wantStatus || err != nil {
				t.Errorf("resize error, expected status %v, having status %v %#v", tt.wantStatus, resp.StatusCode, resp)
			}
		})
	}
}

func Test_handleCrop(t *testing.T) {
	go New().Start(9898)
	type args struct {
		files  []string
		width  string
		height string
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantStatus int
	}{
		{
			name:       "empty",
			args:       args{},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty file",
			args: args{
				width:  "900px",
				height: "900px",
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "empty size",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
			},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "1 file size",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
				width:  "900px",
				height: "900px",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "couple files angle",
			args: args{
				files: []string{
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
				},
				width:  "900px",
				height: "900px",
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			fw, err := w.CreateFormField("width")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(tt.args.width)); err != nil {
				t.Fail()
				return
			}
			fw, err = w.CreateFormField("height")
			if err != nil {
				t.Fail()
				return
			}
			if _, err = fw.Write([]byte(tt.args.height)); err != nil {
				t.Fail()
				return
			}
			for _, file := range tt.args.files {
				f, err := os.Open(file)
				if err != nil {
					t.Fail()
					return
				}
				defer f.Close()
				fw, err = w.CreateFormFile("file", file)
				if err != nil {
					t.Fail()
					return
				}
				if _, err = io.Copy(fw, f); err != nil {
					t.Fail()
					return
				}
			}
			w.Close()
			req, err := http.NewRequest("POST", baseURL+"/crop", &b)
			if err != nil {
				t.Fail()
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			startTime := time.Now().UnixNano()
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano()
			if err != nil {
				t.Fail()
				return
			}
			t.Logf(`request "%v" processing duration: %v sec`, tt.name, float32(endTime-startTime)/1000000000)
			if resp.StatusCode != tt.wantStatus || err != nil {
				t.Errorf("crop error, expected status %v, having status %v", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func Test_handleFlip(t *testing.T) {
	go New().Start(9898)
	type args struct {
		files []string
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantStatus int
	}{
		{
			name:       "empty",
			args:       args{},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "1 file",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "couple files",
			args: args{
				files: []string{
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
				},
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			for _, file := range tt.args.files {
				f, err := os.Open(file)
				if err != nil {
					t.Fail()
					return
				}
				defer f.Close()
				fw, err := w.CreateFormFile("file", file)
				if err != nil {
					t.Fail()
					return
				}
				if _, err = io.Copy(fw, f); err != nil {
					t.Fail()
					return
				}
			}
			w.Close()
			req, err := http.NewRequest("POST", baseURL+"/flip", &b)
			if err != nil {
				t.Fail()
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			startTime := time.Now().UnixNano()
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano()
			if err != nil {
				t.Fail()
				return
			}
			t.Logf(`request "%v" processing duration: %v sec`, tt.name, float32(endTime-startTime)/1000000000)
			if resp.StatusCode != tt.wantStatus || err != nil {
				t.Errorf("flip error, expected status %v, having status %v", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func Test_handleFlop(t *testing.T) {
	go New().Start(9898)
	type args struct {
		files []string
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantStatus int
	}{
		{
			name:       "empty",
			args:       args{},
			want:       model.ErrorResponse{},
			wantStatus: 400,
		},
		{
			name: "1 file",
			args: args{
				files: []string{
					"../../test/image.jpg",
				},
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
		{
			name: "couple files",
			args: args{
				files: []string{
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
					"../../test/image.jpg",
				},
			},
			want:       model.ImagesProcessingResponse{},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			for _, file := range tt.args.files {
				f, err := os.Open(file)
				if err != nil {
					t.Fail()
					return
				}
				defer f.Close()
				fw, err := w.CreateFormFile("file", file)
				if err != nil {
					t.Fail()
					return
				}
				if _, err = io.Copy(fw, f); err != nil {
					t.Fail()
					return
				}
			}
			w.Close()
			req, err := http.NewRequest("POST", baseURL+"/flop", &b)
			if err != nil {
				t.Fail()
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			startTime := time.Now().UnixNano()
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano()
			if err != nil {
				t.Fail()
				return
			}
			t.Logf(`request "%v" processing duration: %v sec`, tt.name, float32(endTime-startTime)/1000000000)
			if resp.StatusCode != tt.wantStatus || err != nil {
				t.Errorf("flop error, expected status %v, having status %v", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
