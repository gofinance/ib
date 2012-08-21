package trade

type Strategy interface {
    Start(engine *Engine)
    Step(message interface{}) bool
    Error(err error)
}
