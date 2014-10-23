package main

import (
  logging "github.com/op/go-logging"
)

// globally available server log
var serverLog *logger = nil

const (
  logFormat = "%{color}%{level} %{module} %{shortfile}%{color:reset} %{message}"
)

// simple wrapper around the go-logging library that allows us to easily create child loggers
// TODO: support filesystem logging
type logger struct {
  *logging.Logger
  level int
}

// get a new server logger
func NewServerLogger(logLevel int, logFormat string) *logger {
  logger := &logger{logging.MustGetLogger("server"), logLevel}

  logging.SetFormatter(logging.MustStringFormatter(logFormat))

  logging.SetLevel(logging.Level(logLevel), "server")
  return logger
}

// get a child logger
func (sl logger) Child(name string) (*logger, error) {
  libLogger, err := logging.GetLogger(name)
  if err != nil {
    return nil, err
  }

  child := &logger{libLogger, sl.level}
  logging.SetLevel(logging.Level(sl.level), name)
  return child, nil
}