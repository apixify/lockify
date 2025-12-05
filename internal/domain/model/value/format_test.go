package value

import "testing"

func TestNewFileFormat(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    FileFormat
		wantErr bool
	}{
		{
			name:    "json format",
			value:   "json",
			want:    Json,
			wantErr: false,
		},
		{
			name:    "dotenv format",
			value:   "dotenv",
			want:    DotEnv,
			wantErr: false,
		},
		{
			name:    "custom format",
			value:   "custom",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "json uppercase",
			value:   "JSON",
			want:    "",
			wantErr: true,
		},
		{
			name:    "dotenv uppercase",
			value:   "DOTENV",
			want:    "",
			wantErr: true,
		},
		{
			name:    "yaml format",
			value:   "yaml",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileFormat(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileFormat(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NewFileFormat(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestFileFormat_String(t *testing.T) {
	tests := []struct {
		name       string
		fileFormat FileFormat
		want       string
	}{
		{
			name:       "json format",
			fileFormat: Json,
			want:       "json",
		},
		{
			name:       "dotenv format",
			fileFormat: DotEnv,
			want:       "dotenv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileFormat.String()
			if got != tt.want {
				t.Errorf("FileFormat.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFileFormat_IsJson(t *testing.T) {
	tests := []struct {
		name       string
		fileFormat FileFormat
		want       bool
	}{
		{
			name:       "json format",
			fileFormat: Json,
			want:       true,
		},
		{
			name:       "dotenv format",
			fileFormat: DotEnv,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileFormat.IsJson()
			if got != tt.want {
				t.Errorf("FileFormat.IsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileFormat_IsDotEnv(t *testing.T) {
	tests := []struct {
		name       string
		fileFormat FileFormat
		want       bool
	}{
		{
			name:       "dotenv format",
			fileFormat: DotEnv,
			want:       true,
		},
		{
			name:       "json format",
			fileFormat: Json,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileFormat.IsDotEnv()
			if got != tt.want {
				t.Errorf("FileFormat.IsDotEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileFormat_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		fileFormat FileFormat
		want       bool
	}{
		{
			name:       "json format",
			fileFormat: Json,
			want:       true,
		},
		{
			name:       "dotenv format",
			fileFormat: DotEnv,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileFormat.IsValid()
			if got != tt.want {
				t.Errorf("FileFormat.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
