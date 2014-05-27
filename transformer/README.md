# transformer
--
    import "github.com/azylman/getl/transformer"


## Usage

#### type Transformer

```go
type Transformer struct {
}
```

A Transformer allows you to easily chain multiple transforms on a table.

#### func  New

```go
func New(table getl.Table) *Transformer
```
New returns a Transformer that allows you to chain transformations on a Table.

#### func (*Transformer) Apply

```go
func (t *Transformer) Apply(transform getl.TransformFunc) *Transformer
```
Apply applies a given TransformFunc to the Transformer.

#### func (*Transformer) Each

```go
func (t *Transformer) Each(transform func(getl.Row) error) *Transformer
```
Each Applies an Each transform.

#### func (*Transformer) Fieldmap

```go
func (t *Transformer) Fieldmap(mappings map[string][]string) *Transformer
```
Fieldmap Applies a Fieldmap transform.

#### func (*Transformer) Map

```go
func (t *Transformer) Map(transform func(getl.Row) (getl.Row, error)) *Transformer
```
Map Applies a Map transform.

#### func (*Transformer) Select

```go
func (t *Transformer) Select(filter func(getl.Row) (bool, error)) *Transformer
```
Select Applies a Select transform.

#### func (Transformer) Table

```go
func (t Transformer) Table() getl.Table
```
Table returns the terminating Table in a Transformer chain.

#### func (*Transformer) TableTransform

```go
func (t *Transformer) TableTransform(transform func(getl.Row, chan<- getl.Row) error) *Transformer
```
TableTransform Applies a TableTransform transform.

#### func (*Transformer) Valuemap

```go
func (t *Transformer) Valuemap(mappings map[string]map[interface{}]interface{}) *Transformer
```
Valuemap Applies a Valuemap transform.