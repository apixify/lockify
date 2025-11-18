package service

import "io"

type ImportService interface {
	FromJson(r io.Reader) (map[string]string, error)
	FromDotEnv(r io.Reader) (map[string]string, error)
}
