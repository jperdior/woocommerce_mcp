package search_products

import (
	"context"
	"fmt"
	"strconv"
	"woocommerce-mcp/internal/product/domain"
)

// ProductSearcher handles product search operations
type ProductSearcher struct {
	productRepository domain.ProductRepository
}

// NewProductSearcher creates a new ProductSearcher
func NewProductSearcher(productRepository domain.ProductRepository) *ProductSearcher {
	return &ProductSearcher{
		productRepository: productRepository,
	}
}

// Execute performs the product search
func (ps *ProductSearcher) Execute(ctx context.Context, request *SearchRequest) (*SearchResponse, error) {
	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Convert request to domain search criteria
	criteria, err := ps.requestToCriteria(request)
	if err != nil {
		return nil, err
	}

	// Validate criteria
	if err := criteria.Validate(); err != nil {
		return nil, err
	}

	// Search products
	products, err := ps.productRepository.Search(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	// Get total count for pagination
	totalCount, err := ps.productRepository.Count(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Convert domain products to response DTOs
	productDTOs := make([]*ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = ps.productToDTO(product)
	}

	// Calculate pagination info
	totalPages := int((totalCount + int64(criteria.PerPage) - 1) / int64(criteria.PerPage))

	return &SearchResponse{
		Products:    productDTOs,
		TotalCount:  int(totalCount),
		CurrentPage: criteria.Page,
		PerPage:     criteria.PerPage,
		TotalPages:  totalPages,
		HasNext:     criteria.Page < totalPages,
		HasPrev:     criteria.Page > 1,
	}, nil
}

// requestToCriteria converts SearchRequest to domain SearchCriteria
func (ps *ProductSearcher) requestToCriteria(request *SearchRequest) (*domain.SearchCriteria, error) {
	criteria := domain.NewSearchCriteria()

	// Set search term
	if request.Search != nil && *request.Search != "" {
		criteria.SetSearch(*request.Search)
	}

	// Set category
	if request.Category != nil && *request.Category != "" {
		criteria.SetCategory(*request.Category)
	}

	// Set tag
	if request.Tag != nil && *request.Tag != "" {
		criteria.SetTag(*request.Tag)
	}

	// Set status
	if request.Status != nil && *request.Status != "" {
		status := domain.ProductStatus(*request.Status)
		if !status.IsValid() {
			return nil, domain.NewInvalidProductStatusError(*request.Status)
		}
		criteria.SetStatus(status)
	}

	// Set type
	if request.Type != nil && *request.Type != "" {
		productType := domain.ProductType(*request.Type)
		if !productType.IsValid() {
			return nil, domain.NewInvalidProductTypeError(*request.Type)
		}
		criteria.SetType(productType)
	}

	// Set featured
	if request.Featured != nil {
		featured, err := strconv.ParseBool(*request.Featured)
		if err != nil {
			return nil, domain.NewProductValidationError("featured", "must be true or false")
		}
		criteria.SetFeatured(featured)
	}

	// Set on sale
	if request.OnSale != nil {
		onSale, err := strconv.ParseBool(*request.OnSale)
		if err != nil {
			return nil, domain.NewProductValidationError("on_sale", "must be true or false")
		}
		criteria.SetOnSale(onSale)
	}

	// Set price range
	var minPrice, maxPrice *domain.Money
	if request.MinPrice != nil && *request.MinPrice != "" {
		price, err := domain.NewMoneyFromString(*request.MinPrice, "USD")
		if err != nil {
			return nil, domain.NewProductValidationError("min_price", "invalid price format")
		}
		minPrice = price
	}
	if request.MaxPrice != nil && *request.MaxPrice != "" {
		price, err := domain.NewMoneyFromString(*request.MaxPrice, "USD")
		if err != nil {
			return nil, domain.NewProductValidationError("max_price", "invalid price format")
		}
		maxPrice = price
	}
	if minPrice != nil || maxPrice != nil {
		criteria.SetPriceRange(minPrice, maxPrice)
	}

	// Set stock status
	if request.StockStatus != nil && *request.StockStatus != "" {
		stockStatus := domain.StockStatus(*request.StockStatus)
		if !stockStatus.IsValid() {
			return nil, domain.NewInvalidStockStatusError(*request.StockStatus)
		}
		criteria.SetStockStatus(stockStatus)
	}

	// Set pagination
	page := 1
	perPage := 10

	if request.Page != nil && *request.Page != "" {
		p, err := strconv.Atoi(*request.Page)
		if err != nil || p < 1 {
			return nil, domain.NewProductValidationError("page", "must be a positive integer")
		}
		page = p
	}

	if request.PerPage != nil && *request.PerPage != "" {
		pp, err := strconv.Atoi(*request.PerPage)
		if err != nil || pp < 1 {
			return nil, domain.NewProductValidationError("per_page", "must be a positive integer")
		}
		perPage = pp
	}

	criteria.SetPagination(page, perPage)

	// Set sorting
	orderBy := "date"
	order := "desc"

	if request.OrderBy != nil && *request.OrderBy != "" {
		orderBy = *request.OrderBy
	}

	if request.Order != nil && *request.Order != "" {
		order = *request.Order
	}

	criteria.SetSorting(orderBy, order)

	return criteria, nil
}

// productToDTO converts domain Product to ProductDTO
func (ps *ProductSearcher) productToDTO(product *domain.Product) *ProductDTO {
	dto := &ProductDTO{
		ID:                product.ID.Value(),
		Name:              product.Name,
		Slug:              product.Slug,
		Permalink:         product.Permalink,
		DateCreated:       product.DateCreated.Format("2006-01-02T15:04:05"),
		DateModified:      product.DateModified.Format("2006-01-02T15:04:05"),
		Type:              string(product.Type),
		Status:            string(product.Status),
		Featured:          product.Featured,
		CatalogVisibility: product.CatalogVisibility,
		Description:       product.Description,
		ShortDescription:  product.ShortDescription,
		SKU:               product.SKU,
		OnSale:            product.OnSale,
		Purchasable:       product.Purchasable,
		TotalSales:        product.TotalSales,
		Virtual:           product.Virtual,
		Downloadable:      product.Downloadable,
		ExternalURL:       product.ExternalURL,
		ButtonText:        product.ButtonText,
		TaxStatus:         product.TaxStatus,
		TaxClass:          product.TaxClass,
		ManageStock:       product.ManageStock,
		StockQuantity:     product.StockQuantity,
		StockStatus:       string(product.StockStatus),
		Backorders:        product.Backorders,
		BackordersAllowed: product.BackordersAllowed,
		Backordered:       product.Backordered,
		Weight:            product.Weight,
		ShippingRequired:  product.ShippingRequired,
		ShippingTaxable:   product.ShippingTaxable,
		ShippingClass:     product.ShippingClass,
		ShippingClassID:   product.ShippingClassID,
		ReviewsAllowed:    product.ReviewsAllowed,
		AverageRating:     product.AverageRating,
		RatingCount:       product.RatingCount,
		RelatedIDs:        product.RelatedIDs,
		UpsellIDs:         product.UpsellIDs,
		CrossSellIDs:      product.CrossSellIDs,
		ParentID:          product.ParentID,
		PurchaseNote:      product.PurchaseNote,
		Variations:        product.Variations,
		GroupedProducts:   product.GroupedProducts,
		MenuOrder:         product.MenuOrder,
	}

	// Convert price
	if product.Price != nil {
		priceStr := fmt.Sprintf("%.2f", product.Price.Amount())
		dto.Price = priceStr
	}

	// Convert regular price
	if product.RegularPrice != nil {
		regularPriceStr := fmt.Sprintf("%.2f", product.RegularPrice.Amount())
		dto.RegularPrice = regularPriceStr
	}

	// Convert sale price
	if product.SalePrice != nil {
		salePriceStr := fmt.Sprintf("%.2f", product.SalePrice.Amount())
		dto.SalePrice = salePriceStr
	}

	// Convert dimensions
	if product.Dimensions != nil {
		dto.Dimensions = &DimensionsDTO{
			Length: product.Dimensions.Length,
			Width:  product.Dimensions.Width,
			Height: product.Dimensions.Height,
		}
	}

	// Convert categories
	dto.Categories = make([]*CategoryDTO, len(product.Categories))
	for i, category := range product.Categories {
		dto.Categories[i] = &CategoryDTO{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		}
	}

	// Convert tags
	dto.Tags = make([]*TagDTO, len(product.Tags))
	for i, tag := range product.Tags {
		dto.Tags[i] = &TagDTO{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	// Convert images
	dto.Images = make([]*ImageDTO, len(product.Images))
	for i, image := range product.Images {
		dto.Images[i] = &ImageDTO{
			ID:           image.ID,
			DateCreated:  image.DateCreated,
			DateModified: image.DateModified,
			Src:          image.Src,
			Name:         image.Name,
			Alt:          image.Alt,
			Position:     image.Position,
		}
	}

	// Convert attributes
	dto.Attributes = make([]*AttributeDTO, len(product.Attributes))
	for i, attribute := range product.Attributes {
		dto.Attributes[i] = &AttributeDTO{
			ID:        attribute.ID,
			Name:      attribute.Name,
			Position:  attribute.Position,
			Visible:   attribute.Visible,
			Variation: attribute.Variation,
			Options:   attribute.Options,
		}
	}

	// Convert default attributes
	dto.DefaultAttributes = make([]*DefaultAttributeDTO, len(product.DefaultAttributes))
	for i, defaultAttr := range product.DefaultAttributes {
		dto.DefaultAttributes[i] = &DefaultAttributeDTO{
			ID:     defaultAttr.ID,
			Name:   defaultAttr.Name,
			Option: defaultAttr.Option,
		}
	}

	// Convert metadata
	dto.MetaData = make([]*MetaDataDTO, len(product.MetaData))
	for i, metaData := range product.MetaData {
		dto.MetaData[i] = &MetaDataDTO{
			ID:    metaData.ID,
			Key:   metaData.Key,
			Value: metaData.Value,
		}
	}

	return dto
}
