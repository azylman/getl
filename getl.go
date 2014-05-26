package getl

// Table is a representation of a table of data.
type Table interface {
	// Rows returns a channel that provides the Rows in the table.
	Rows() <-chan Row
	// Err returns the first non-EOF error that was encountered by the Table.
	Err() error
	// Stop signifies that a Table should stop sending Rows down its channel.
	// A Table is also responsible for calling Stop on any upstream Tables it knows about.
	// Stop should be idempotent. It's expected that Stop will never be called by a consumer of a
	// Table unless that consumer is also a Table. It can be used to Stop all upstream Tables in
	// the event of an error that needs to halt the pipeline.
	Stop()
}

// Row is a representation of a line of data in a Table.
type Row map[string]interface{}

// TransformFunc is a function that can be applied to a Table to transform it. It should receive the
// Rows from in and may send any number of Rows to out. It should not return until it has finished
// all work (received all the Rows it's going to receive, sent all the Rows it's going to send).
type TransformFunc func(in <-chan Row, out chan<- Row) error

// Transform returns a new Table that provides all the Rows of the input Table transformed with the TransformFunc.
func Transform(source Table, transform TransformFunc) Table {
	return newTransformedTable(source, transform)
}

type transformedTable struct {
	source  Table
	err     error
	rows    chan Row
	stopped bool
}

func (t transformedTable) Rows() <-chan Row {
	return t.rows
}

func (t transformedTable) Err() error {
	return t.err
}

func (t *transformedTable) Stop() {
	if t.stopped {
		return
	}
	t.stopped = true
	t.source.Stop()
}

func drain(c <-chan Row) {
	for _ = range c {
		// Drain everything left in the channel
	}
}

func (t *transformedTable) start(transform TransformFunc) {
	// A level of indirection is necessary between the i/o channels and the TransformFunc so that
	// the TransformFunc doesn't need to know about the stop state of any of the Tables.
	in := make(chan Row)
	out := make(chan Row)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	stop := func() {
		if t.stopped {
			return
		}
		t.Stop()
		drain(t.source.Rows())
		drain(out)
		close(t.rows)
	}
	defer stop()

	go func() {
		defer close(errChan)
		defer close(out)
		if err := transform(in, out); err != nil {
			errChan <- err
		}
	}()
	go func() {
		defer func() {
			doneChan <- struct{}{}
		}()
		for row := range out {
			if t.stopped {
				continue
			}
			t.rows <- row
		}
	}()

	go func() {
		defer func() {
			doneChan <- struct{}{}
		}()
		defer close(in)
		for row := range t.source.Rows() {
			if t.stopped {
				continue
			}
			in <- row
		}
	}()
	for err := range errChan {
		stop()
		t.err = err
	}
	// Wait for all channels to finish
	<-doneChan // Once to make sure we've consumed the output of the TransformFunc
	<-doneChan // Once to make sure we've consumed the output of the source Table
	if t.source.Err() != nil {
		t.err = t.source.Err()
	}
}

func newTransformedTable(source Table, transform TransformFunc) Table {
	table := &transformedTable{
		source: source,
		rows:   make(chan Row),
	}
	go table.start(transform)
	return table
}
