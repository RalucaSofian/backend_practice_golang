package query_utils

import (
	"app/models"
	"app/utils"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
)

// Possible Filter Operations
type FilterOperation string

const (
	FilterOp_EQ      FilterOperation = "EQ"
	FilterOp_IS_NULL FilterOperation = "IS_NULL"
	FilterOp_IN      FilterOperation = "IN"
	FilterOp_NOT_IN  FilterOperation = "NOT_IN"
	FilterOp_GT      FilterOperation = "GT"
	FilterOp_GTE     FilterOperation = "GTE"
	FilterOp_LT      FilterOperation = "LT"
	FilterOp_LTE     FilterOperation = "LTE"
)

// Filter Info structure
type FilterInfo struct {
	Field     string
	Operation FilterOperation
	Value     interface{}
}

// Order Direction (ASC/DESC)
type OrderDirection string

const (
	OrderDir_ASC  OrderDirection = "ASC"
	OrderDir_DESC OrderDirection = "DESC"
)

// Order Info structure
type OrderInfo struct {
	Field     string
	Direction OrderDirection
}

// Query Info object structure
type QueryInfo struct {
	Search  *string
	Filters *[]FilterInfo
	Order   *OrderInfo
	Limit   *int
	Offset  *int
}

// Modify a base Bun Query, according to the Query Info object
func (queryInfo *QueryInfo) Process(baseQuery *bun.SelectQuery) *bun.SelectQuery {
	fmt.Println("[query] Process Query")
	query := baseQuery

	if queryInfo.Search != nil {
		searchTerm := *queryInfo.Search
		searchTerm = strings.ReplaceAll(searchTerm, " ", " & ")
		query = query.Where("?TableAlias.search_vector @@ to_tsquery('simple', ? || ':*')", searchTerm)
	}

	if queryInfo.Order != nil {
		orderStr := fmt.Sprintf("%s %s", queryInfo.Order.Field, queryInfo.Order.Direction)
		query = query.Order(orderStr)
	}

	if queryInfo.Limit != nil {
		query = query.Limit(*queryInfo.Limit)
	}

	if queryInfo.Offset != nil {
		query = query.Offset(*queryInfo.Offset)
	}

	if queryInfo.Filters != nil {
		for _, filter := range *queryInfo.Filters {
			switch filter.Operation {
			case FilterOp_EQ:
				query = query.Where("? = ?", bun.Ident(filter.Field), filter.Value)

			case FilterOp_IS_NULL:
				if filter.Value == true {
					query = query.Where("? IS NULL", bun.Ident(filter.Field))
				} else if filter.Value == false {
					query = query.Where("? IS NOT NULL", bun.Ident(filter.Field))
				}

			case FilterOp_IN:
				query = query.Where("? IN (?)", bun.Ident(filter.Field), bun.In(filter.Value))

			case FilterOp_NOT_IN:
				query = query.Where("? NOT IN (?)", bun.Ident(filter.Field), bun.In(filter.Value))

			case FilterOp_GT:
				query = query.Where("? > ?", bun.Ident(filter.Field), filter.Value)

			case FilterOp_GTE:
				query = query.Where("? >= ?", bun.Ident(filter.Field), filter.Value)

			case FilterOp_LT:
				query = query.Where("? < ?", bun.Ident(filter.Field), filter.Value)

			case FilterOp_LTE:
				query = query.Where("? <= ?", bun.Ident(filter.Field), filter.Value)

			default:
				query = query.Where("? = ?", bun.Ident(filter.Field), filter.Value)
			}
		}
	}

	fmt.Printf("[query] %s\n", query)
	return query
}

func strToOrderDir(inputStr string) OrderDirection {
	input := strings.ToLower(inputStr)
	switch input {
	case "asc":
		return OrderDir_ASC
	case "desc", "dsc":
		return OrderDir_DESC
	default:
		return OrderDir_ASC
	}
}

func strToFilterOper(inputStr string) FilterOperation {
	input := strings.ToLower(inputStr)
	switch input {
	case "eq":
		return FilterOp_EQ
	case "is_null":
		return FilterOp_IS_NULL
	case "in":
		return FilterOp_IN
	case "not_in":
		return FilterOp_NOT_IN
	case "gt":
		return FilterOp_GT
	case "gte":
		return FilterOp_GTE
	case "lt":
		return FilterOp_LT
	case "lte":
		return FilterOp_LTE
	default:
		return FilterOp_EQ
	}
}

// Check if the Field from Query Info is part of the Model
func validateFieldName[T models.Model](fieldName string) bool {
	var modelInstance T
	model := reflect.TypeOf(modelInstance)

	for i := range model.NumField() {
		modelField := model.Field(i)

		bunTag := modelField.Tag.Get("bun")
		tagsList := strings.Split(bunTag, ",")
		nameTag := tagsList[0]

		isTableName := strings.HasPrefix(nameTag, "table:")
		if !isTableName {
			if nameTag == fieldName {
				fmt.Println("[query] Model", model.Name(), "contains field:", fieldName)
				return true
			}
		}
	}

	fmt.Println("[query] Model", model.Name(), "does not contain field:", fieldName)
	return false
}

// Parse the Query Params of a HTTP Req into a Query Object
func ParseQueryParams[T models.Model](input map[string][]string) (*QueryInfo, error) {
	queryObj := QueryInfo{}
	filtersList := []FilterInfo{}

	if len(input) == 0 {
		fmt.Println("[query] No Query present")
		return nil, nil
	}

	for key, val := range input {
		key = strings.ToLower(key)

		if key == "search" {
			search := strings.ToLower(val[0])
			queryObj.Search = &search

		} else if key == "limit" {
			limitAsInt, _ := strconv.Atoi(val[0])
			queryObj.Limit = &limitAsInt

		} else if key == "offset" {
			offsetAsInt, _ := strconv.Atoi(val[0])
			queryObj.Offset = &offsetAsInt

		} else if key == "order" {
			fieldAndDir := strings.Split(val[0], "__")
			if len(fieldAndDir) == 1 { // only Field (default ASC)
				fmt.Println("[query] Default Order Direction: Ascending")
				orderField := fieldAndDir[0]
				if !validateFieldName[T](orderField) {
					return nil, utils.NewApiError(utils.ErrorType_QueryError,
						fmt.Sprintf("Invalid Field Name: %s", orderField))
				}
				queryObj.Order = &OrderInfo{
					Field:     orderField,
					Direction: OrderDir_ASC,
				}
			} else if len(fieldAndDir) == 2 { // Field AND Direction
				orderField := fieldAndDir[0]
				if !validateFieldName[T](orderField) {
					return nil, utils.NewApiError(utils.ErrorType_QueryError,
						fmt.Sprintf("Invalid Field Name: %s", orderField))
				}
				orderDir := fieldAndDir[1]
				direction := strToOrderDir(orderDir)
				queryObj.Order = &OrderInfo{
					Field:     orderField,
					Direction: direction,
				}
			} else {
				fmt.Println("[query] Invalid number of arguments in Order Query")
				return nil, utils.NewApiError(utils.ErrorType_QueryError, "Invalid number of arguments")
			}

		} else { // Filters
			var filterObj FilterInfo
			fieldAndOp := strings.Split(key, "__")

			if len(fieldAndOp) == 1 { // only Field (default Op: EQ)
				fmt.Println("[query] Default Filter Operation: Equal")
				filterField := fieldAndOp[0]
				if !validateFieldName[T](filterField) {
					return nil, utils.NewApiError(utils.ErrorType_QueryError,
						fmt.Sprintf("Invalid Field Name: %s", filterField))
				}
				filterVal := val[0]
				filterObj = FilterInfo{
					Field:     filterField,
					Operation: FilterOp_EQ,
					Value:     filterVal,
				}
				filtersList = append(filtersList, filterObj)

			} else if len(fieldAndOp) == 2 { // Field AND Operation
				filterField := fieldAndOp[0]
				//
				if !validateFieldName[T](filterField) {
					return nil, utils.NewApiError(utils.ErrorType_QueryError,
						fmt.Sprintf("Invalid Field Name: %s", filterField))
				}
				filterOp := fieldAndOp[1]
				operation := strToFilterOper(filterOp)
				filterVal := val[0]

				if operation == FilterOp_IS_NULL { // Filter Op: Is_Null
					if strings.ToLower(filterVal) == "true" { // Is_Null=TRUE
						filterObj = FilterInfo{
							Field:     filterField,
							Operation: operation,
							Value:     true,
						}
					} else if strings.ToLower(filterVal) == "false" { // Is_Null=FALSE
						filterObj = FilterInfo{
							Field:     filterField,
							Operation: operation,
							Value:     false,
						}
					} else {
						fmt.Println("[query] Invalid Filter Value. Expected true/false for", key)
						return nil, utils.NewApiError(utils.ErrorType_QueryError,
							fmt.Sprintf("Invalid Filter Value. Expected true/false for %s", key))
					}
				} else if operation == FilterOp_IN || operation == FilterOp_NOT_IN { // Filter Op: In/Not_In=LIST
					listOfVals := strings.Split(filterVal, ",")
					filterObj = FilterInfo{
						Field:     filterField,
						Operation: operation,
						Value:     listOfVals,
					}
				} else { // rest of Filter Operations
					filterObj = FilterInfo{
						Field:     filterField,
						Operation: operation,
						Value:     filterVal,
					}
				}
				filtersList = append(filtersList, filterObj)
			} else {
				fmt.Println("[query] Invalid number of arguments in Filter Query")
				return nil, utils.NewApiError(utils.ErrorType_QueryError, "Invalid number of arguments")
			}
			queryObj.Filters = &filtersList
		}
	}
	fmt.Println("[query] -- QUERY --")
	if queryObj.Search != nil {
		fmt.Println("[query] SEARCH =", *queryObj.Search)
	}
	if queryObj.Filters != nil {
		fmt.Println("[query] FILTERS =", *queryObj.Filters)
	}
	if queryObj.Order != nil {
		fmt.Println("[query] ORDER =", *queryObj.Order)
	}
	if queryObj.Limit != nil {
		fmt.Println("[query] LIMIT =", *queryObj.Limit)
	}
	if queryObj.Offset != nil {
		fmt.Println("[query] OFFSET =", *queryObj.Offset)
	}
	fmt.Println("[query] -- --")

	return &queryObj, nil
}
