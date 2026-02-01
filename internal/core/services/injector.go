package services

import "go.uber.org/dig"

func NewInjector(digger *dig.Container) {
	digger.Provide(NewDeviceMessageService)
}
