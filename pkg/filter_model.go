package pkg

type FilterModel struct {
	SearchString *string `json:"search_string,omitempty"`
	ProductType  *string `json:"product_type,omitempty"`
	ProductBrand *string `json:"product_brand,omitempty"`
	ProductName  *string `json:"product_name,omitempty"`
	ProductModel *string `json:"product_model,omitempty"`
	MinStock     *int    `json:"min_stock,omitempty"`
	MaxStock     *int    `json:"max_stock,omitempty"`
	Supplier     *string `json:"supplier,omitempty"`
}
