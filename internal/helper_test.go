package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitShellCommand(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{name: "1", args: `a b c`, want: []string{"a", "b", "c"}},
		{name: "2", args: `a b "c d e"`, want: []string{"a", "b", "c d e"}},
		{name: "3", args: `a b " c d e "`, want: []string{"a", "b", " c d e "}},
		{name: "4", args: `  a   b   c `, want: []string{"a", "b", "c"}},
		{name: "5", args: `a b 'c'`, want: []string{"a", "b", "c"}},
		{name: "6", args: `a "b" 'c'`, want: []string{"a", "b", "c"}},
		{name: "7", args: `a "b " 'c '`, want: []string{"a", "b ", "c "}},
		{name: "8", args: `\`, want: []string{`\`}},
		{name: "9", args: "\\", want: []string{`\`}},
		{name: "10", args: `\ ""`, want: []string{`\`, ""}},
		{name: "11", args: `\ "" "c "`, want: []string{`\`, "", "c "}},
		{name: "12", args: `\ "" "c\ "`, want: []string{`\`, "", "c\\ "}},
		{name: "12", args: "'a' 'b c' '\"'", want: []string{`a`, `b c`, `"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SplitShellCommand(tt.args), tt.name+" / "+tt.args)
		})
	}
}

func TestGetRemoteRevPath(t *testing.T) {
	type args struct {
		localRootPath     string
		remoteRootPath    string
		localPath         string
		keepRemoteBaseDir bool
	}
	tests := []struct {
		name       string
		args       args
		want       string
		errContain string
	}{
		{
			name: "1",
			args: args{localRootPath: "./", remoteRootPath: "/root", localPath: "a.txt", keepRemoteBaseDir: false},
			want: "/root/a.txt",
		},
		{
			name: "2",
			args: args{localRootPath: ".", remoteRootPath: "/root", localPath: "a.txt", keepRemoteBaseDir: false},
			want: "/root/a.txt",
		},
		{
			name: "3",
			args: args{localRootPath: "./dir", remoteRootPath: "/root", localPath: "dir/a.txt", keepRemoteBaseDir: false},
			want: "/root/a.txt",
		},

		{
			name: "3",
			args: args{localRootPath: "./dir", remoteRootPath: "/root", localPath: "dir/root/a.txt", keepRemoteBaseDir: true},
			want: "/root/a.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, GetRemoteRevPath(tt.args.localRootPath, tt.args.remoteRootPath, tt.args.localPath, tt.args.keepRemoteBaseDir), tt.name+" / "+tt.args.localRootPath)
			})
		})
	}
}
