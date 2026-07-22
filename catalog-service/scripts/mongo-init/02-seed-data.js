db = db.getSiblingDB("catalog_db");

const electronicsId = new ObjectId();
const audioSubCatId = new ObjectId();
const clothingId = new ObjectId();

db.categories.insertMany([
  {
    _id: electronicsId,
    parent_id: null,
    name_vi: "Thiết bị điện tử",
    name_en: "Electronics",
    slug: "electronics",
    sort_order: 1,
    created_at: new Date(),
  },
  {
    _id: audioSubCatId,
    parent_id: electronicsId,
    name_vi: "Âm thanh",
    name_en: "Audio",
    slug: "audio",
    sort_order: 2,
    created_at: new Date(),
  },
  {
    _id: clothingId,
    parent_id: null,
    name_vi: "Quần áo",
    name_en: "Clothing",
    slug: "clothing",
    sort_order: 1,
    created_at: new Date(),
  },
]);

db.products.insertMany([
  {
    category_id: audioSubCatId,
    sku: "HEADPHONE-001",
    name_vi: "Tai nghe chống ồn không dây",
    name_en: "Wireless Noise Canceling Headphones",
    description_vi: "Mô tả sản phẩm...",
    description_en: "Product description...",
    unit: "piece",
    base_price: NumberLong(200000000), // 2,000,000 VND scaled by 10^2 or $200.00 scaled by 10^5
    sale_price: NumberLong(180000000),
    rating_avg: 4.8,
    rating_count: 150,
    is_active: true,
    variants: [
      {
        id: "var-black",
        variant_label: "Color: Black",
        price_delta: NumberLong(0),
        sku: "HEADPHONE-001-BLK",
      },
      {
        id: "var-silver",
        variant_label: "Color: Silver",
        price_delta: NumberLong(10000000), // +$10.00
        sku: "HEADPHONE-001-SLV",
      },
    ],
    images: [
      {
        id: "img-1",
        url: "https://cdn.example.com/products/headphone-1.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("=== Catalog Database successfully seeded ===");
