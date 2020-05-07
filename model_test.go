package exceptionist

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"text/template"
)

func Test_NewBucket(t *testing.T) {
	//when
	bucket := newBucket(TR)

	//then
	assert.Equal(t, 1, len(bucket.rows))
	assert.Contains(t, bucket.rows, "default", "Bucket should be initialized with the `default` row")
	assert.Equal(t, 1, len(bucket.messageTemplates.Templates()))
}

func TestBucket_AddRow(t *testing.T) {
	//given
	var bucket = prepareEmptyTestBucket()

	//when
	bucket.addRow("newRow", 5000, "template")

	//then
	assert.Contains(t, bucket.rows, "newRow", "Row with key`newRow` should be added")
	assert.Equal(t, 5000, bucket.rows["newRow"].errorCode)
	assert.Equal(t, "newRow", bucket.rows["newRow"].templateName)

	template := bucket.messageTemplates.Templates()[1]
	assert.Contains(t, template.Name(), "newRow", "Message Template with name`newRow` should be added")
	assert.Equal(t, reflect.ValueOf(template).Elem().FieldByName("text").String(), "template")
}

func TestBucket_AddRow_Should_Rewrite_When_SameKeyOccurs(t *testing.T) {
	//given
	var bucket = prepareEmptyTestBucket()

	bucket.addRow("newRow", 5000, "template")

	//when
	bucket.addRow("newRow", 5001, "template2")

	//then
	assert.Contains(t, bucket.rows, "newRow", "Row with key`newRow` should be added")
	assert.Equal(t, 5001, bucket.rows["newRow"].errorCode)
	assert.Equal(t, "newRow", bucket.rows["newRow"].templateName)

	tmpl := bucket.messageTemplates.Templates()[1]
	assert.Contains(t, tmpl.Name(), "newRow", "Message Template with name`newRow` should be added")
	assert.Equal(t, "template2", readTemplateValue(tmpl))
}

func TestBucket_FindRow_Return_When_KeyExists(t *testing.T) {
	//given
	var bucket = prepareEmptyTestBucket()

	bucket.addRow("default", 100, "defaultTemplate")
	bucket.addRow("row1", 5000, "template1")

	//when
	row := bucket.findRow("row1")

	//then
	assert.NotNil(t, row)
	assert.Equal(t, 5000, row.errorCode)
	assert.Equal(t, "row1", row.templateName)
}

func TestBucket_FindRow_Return_DefaultRow_When_KeyNotExists(t *testing.T) {
	//given
	var bucket = prepareEmptyTestBucket()

	bucket.addRow("default", 100, "defaultTemplate")

	//when
	row := bucket.findRow("row1")

	//then
	assert.NotNil(t, row)
	assert.Equal(t, 100, row.errorCode)
	assert.Equal(t, "default", row.templateName)
}

func TestBucket_FormatToErrorMessage_Without_Args(t *testing.T) {
	//given
	rows := make(map[string]row)
	row := prepareRow("template1")
	rows["template1"] = row
	tmpl := template.Must(template.New("template1").Parse("template content"))

	bucket := prepareTestBucket(rows, tmpl)

	//when
	errorMessage := bucket.formatToErrorMessage(row, nil)

	//then
	assert.Equal(t, "template content", errorMessage)
}

func TestBucket_FormatToErrorMessage_With_Args(t *testing.T) {
	//given
	rows := make(map[string]row)
	row := prepareRow("template1")
	rows["template1"] = row
	tmpl := template.Must(template.New("template1").Parse("template content arg1: {{index . 0}} arg2: {{index . 1}}"))

	bucket := prepareTestBucket(rows, tmpl)

	//when
	errorMessage := bucket.formatToErrorMessage(row, []interface{}{"1111",  "ğ ş İ ö ç : ; !"})

	//then
	assert.Equal(t, "template content arg1: 1111 arg2: ğ ş İ ö ç : ; !", errorMessage)
}

func Test_NewConfig(t *testing.T) {
	//given
	dir := "myDir"
	prefix := "myPrefix"

	//when
	config := NewConfig(&dir, &prefix)

	//then
	assert.Equal(t, "myDir", *config.dir)
	assert.Equal(t, "myPrefix", *config.prefix)
}

func Test_NewConfig_When_PrefixIsNotGiven(t *testing.T) {
	//given
	dir := "myDir"

	//when
	config := NewConfig(&dir, nil)

	//then
	assert.Equal(t, "myDir", *config.dir)
	assert.Equal(t, "messages", *config.prefix)
}

func prepareEmptyTestBucket() bucket {
	tmpl := template.Must(template.New("test").Parse("test"))
	return prepareTestBucket(map[string]row{}, tmpl)
}

func prepareTestBucket(rows map[string]row, tmpl *template.Template) bucket {
	return bucket{
		rows:             rows,
		messageTemplates: tmpl,
	}
}

func readTemplateValue(tmpl *template.Template) string {
	return reflect.ValueOf(tmpl).Elem().FieldByName("text").String()
}

func prepareRow(templateName string) row {
	return row{templateName: templateName}
}
