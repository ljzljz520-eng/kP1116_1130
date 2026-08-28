package particles

type Emitter struct {
	field Field
	form  string
}

func NewEmitter(count int) *Emitter {
	return &Emitter{field: NewField(count, 1920, 1080), form: "drift"}
}

func (e *Emitter) Advance(step int) Snapshot {
	e.field.Drift(step)
	return e.field.Snapshot(e.form)
}

func (e *Emitter) SetForm(form string) Snapshot {
	if form == "heart" {
		e.form = form
	} else {
		e.form = "drift"
	}
	return e.field.Snapshot(e.form)
}

func (e *Emitter) Current() Snapshot {
	return e.field.Snapshot(e.form)
}

func (e *Emitter) Reset() Snapshot {
	e.form = "drift"
	e.field.Frame = 0
	return e.Current()
}
