package dynamodb_test

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	//
	// The provider linter wants your imports to be in two groups: first,
	// standard library (i.e., "fmt" or "strings"), second, everything else.
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDynamoDBTableItemDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_dynamodb_table_item.test"
	hashKey := "hashKey"
	itemContent := `{
	"hashKey": {"S": "something"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`
	key := `{
	"hashKey": {"S": "something"}
}`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(dynamodb.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemDataSourceConfig_basic(rName, hashKey, itemContent, key),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrEquivalentJSON(dataSourceName, "item", itemContent),
					resource.TestCheckResourceAttr(dataSourceName, "table_name", rName),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItemDataSource_projectionExpression(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_dynamodb_table_item.test"
	hashKey := "hashKey"
	projectionExpression := "one,two"
	itemContent := `{
	"hashKey": {"S": "something"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`
	key := `{
	"hashKey": {"S": "something"}
}`

	expected := `{
	"one": {"N": "11111"},
	"two": {"N": "22222"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(dynamodb.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemDataSourceConfig_ProjectionExpression(rName, hashKey, itemContent, projectionExpression, key),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrEquivalentJSON(dataSourceName, "item", expected),
					resource.TestCheckResourceAttr(dataSourceName, "table_name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "projection_expression", projectionExpression),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItemDataSource_expressionAttributeNames(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_dynamodb_table_item.test"
	hashKey := "hashKey"
	itemContent := `{
	"hashKey": {"S": "something"},
	"one": {"N": "11111"},
	"Percentile": {"N": "22222"}
}`
	key := `{
	"hashKey": {"S": "something"}
}`

	expected := `{
	"Percentile": {"N": "22222"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(dynamodb.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemDataSourceConfig_ExpressionAttributeNames(rName, hashKey, itemContent, key),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrEquivalentJSON(dataSourceName, "item", expected),
					resource.TestCheckResourceAttr(dataSourceName, "table_name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "projection_expression", "#P"),
				),
			},
		},
	})
}

func testAccTableItemDataSourceConfig_basic(tableName, hashKey, item string, key string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key

  item = <<ITEM
%s
ITEM
}

data "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name

  key = <<KEY
%s
KEY
  depends_on = [aws_dynamodb_table_item.test]
}
`, tableName, hashKey, hashKey, item, key)
}

func testAccTableItemDataSourceConfig_ProjectionExpression(tableName, hashKey, item, projectionExpression, key string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key

  item = <<ITEM
%s
ITEM
}

data "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  projection_expression = "%s"
  key = <<KEY
%s
KEY
  depends_on = [aws_dynamodb_table_item.test]
}
`, tableName, hashKey, hashKey, item, projectionExpression, key)
}

func testAccTableItemDataSourceConfig_ExpressionAttributeNames(tableName, hashKey, item string, key string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key

  item = <<ITEM
%s
ITEM
}

data "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  expression_attribute_names = {
	"#P" = "Percentile"
}
  projection_expression = "#P"
  key = <<KEY
%s
KEY
  depends_on = [aws_dynamodb_table_item.test]
}
`, tableName, hashKey, hashKey, item, key)
}
