package module



type AppModuleAdapter interface {
	AppModule
	// RegisterServices allows a module to register services
	RegisterServices(Configurator)
}
