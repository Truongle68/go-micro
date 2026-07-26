db = db.getSiblingDB("catalog-db");

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
  // 1. Audio Product
  {
    category_id: audioSubCatId,
    sku: "HEADPHONE-001",
    name_vi: "Tai nghe chống ồn không dây",
    name_en: "Wireless Noise Canceling Headphones",
    description_vi: "Tai nghe cao cấp tích hợp công nghệ chống ồn chủ động ANC.",
    description_en: "Premium headphones with active noise cancellation technology.",
    unit: "piece",
    base_price: NumberLong(200000000),
    sale_price: NumberLong(180000000),
    rating_avg: 4.8,
    rating_count: 150,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Color: Black",
        price_delta: NumberLong(0),
        sku: "HEADPHONE-001-BLK",
      },
      {
        id: new ObjectId(),
        variant_label: "Color: Silver",
        price_delta: NumberLong(10000000),
        sku: "HEADPHONE-001-SLV",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/headphone-1.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 2. Audio Product
  {
    category_id: audioSubCatId,
    sku: "SPEAKER-002",
    name_vi: "Loa Bluetooth di động kháng nước",
    name_en: "Waterproof Portable Bluetooth Speaker",
    description_vi: "Loa di động âm thanh nổi bass sâu, chuẩn kháng nước IPX7.",
    description_en: "Portable stereo speaker with deep bass and IPX7 waterproof rating.",
    unit: "piece",
    base_price: NumberLong(120000000),
    sale_price: NumberLong(99000000),
    rating_avg: 4.6,
    rating_count: 85,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Color: Blue",
        price_delta: NumberLong(0),
        sku: "SPEAKER-002-BLU",
      },
      {
        id: new ObjectId(),
        variant_label: "Color: Red",
        price_delta: NumberLong(0),
        sku: "SPEAKER-002-RED",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/speaker-2.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 3. Audio Product
  {
    category_id: audioSubCatId,
    sku: "EARBUDS-003",
    name_vi: "Tai nghe In-Ear True Wireless",
    name_en: "True Wireless In-Ear Earbuds",
    description_vi: "Tai nghe nhét tai nhỏ gọn, thời lượng pin lên đến 24 giờ cùng hộp sạc.",
    description_en: "Compact earbuds with up to 24 hours total battery life with charging case.",
    unit: "piece",
    base_price: NumberLong(150000000),
    sale_price: NumberLong(135000000),
    rating_avg: 4.5,
    rating_count: 210,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Edition: Standard",
        price_delta: NumberLong(0),
        sku: "EARBUDS-003-STD",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/earbuds-3.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 4. Electronics Product
  {
    category_id: electronicsId,
    sku: "WATCH-004",
    name_vi: "Đồng hồ thông minh theo dõi sức khỏe",
    name_en: "Smart Fitness Tracker Watch",
    description_vi: "Màn hình AMOLED, đo nhịp tim liên tục và theo dõi giấc ngủ.",
    description_en: "AMOLED screen, continuous heart rate monitoring, and sleep tracking.",
    unit: "piece",
    base_price: NumberLong(250000000),
    sale_price: NumberLong(220000000),
    rating_avg: 4.7,
    rating_count: 310,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Size: 40mm",
        price_delta: NumberLong(0),
        sku: "WATCH-004-40",
      },
      {
        id: new ObjectId(),
        variant_label: "Size: 44mm",
        price_delta: NumberLong(20000000),
        sku: "WATCH-004-44",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/watch-4.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 5. Electronics Product
  {
    category_id: electronicsId,
    sku: "CHARGER-005",
    name_vi: "Sạc nhanh GaN 65W 3 cổng",
    name_en: "65W 3-Port GaN Fast Charger",
    description_vi: "Công nghệ GaN III nhỏ gọn, hỗ trợ sạc nhanh cho Laptop và Điện thoại.",
    description_en: "Compact GaN III technology supporting fast charge for Laptops and Phones.",
    unit: "piece",
    base_price: NumberLong(60000000),
    sale_price: NumberLong(49000000),
    rating_avg: 4.9,
    rating_count: 520,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Plug: US Standard",
        price_delta: NumberLong(0),
        sku: "CHARGER-005-US",
      },
      {
        id: new ObjectId(),
        variant_label: "Plug: EU Standard",
        price_delta: NumberLong(0),
        sku: "CHARGER-005-EU",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/charger-5.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 6. Clothing Product
  {
    category_id: clothingId,
    sku: "TSHIRT-006",
    name_vi: "Áo phông nam Cotton Oversized",
    name_en: "Men's Oversized Cotton T-Shirt",
    description_vi: "Chất liệu 100% cotton thoáng mát, phom dáng rộng thời trang.",
    description_en: "100% breathable cotton fabric with a stylish relaxed fit.",
    unit: "piece",
    base_price: NumberLong(30000000),
    sale_price: NumberLong(25000000),
    rating_avg: 4.4,
    rating_count: 95,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Size M / White",
        price_delta: NumberLong(0),
        sku: "TSHIRT-006-M-WHT",
      },
      {
        id: new ObjectId(),
        variant_label: "Size L / White",
        price_delta: NumberLong(0),
        sku: "TSHIRT-006-L-WHT",
      },
      {
        id: new ObjectId(),
        variant_label: "Size XL / Black",
        price_delta: NumberLong(0),
        sku: "TSHIRT-006-XL-BLK",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/tshirt-6.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 7. Clothing Product
  {
    category_id: clothingId,
    sku: "HOODIE-007",
    name_vi: "Áo Hoodie nỉ dầy Unisex",
    name_en: "Unisex Heavyweight Fleece Hoodie",
    description_vi: "Áo hoodie chất vải nỉ bông giữ ấm tốt, túi kangaroo tiện lợi.",
    description_en: "Warm fleece material hoodie with a comfortable front kangaroo pocket.",
    unit: "piece",
    base_price: NumberLong(55000000),
    sale_price: NumberLong(45000000),
    rating_avg: 4.7,
    rating_count: 140,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Size M / Grey",
        price_delta: NumberLong(0),
        sku: "HOODIE-007-M-GRY",
      },
      {
        id: new ObjectId(),
        variant_label: "Size L / Navy",
        price_delta: NumberLong(0),
        sku: "HOODIE-007-L-NVY",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/hoodie-7.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 8. Clothing Product
  {
    category_id: clothingId,
    sku: "JEANS-008",
    name_vi: "Quần Jeans nam dáng đứng",
    name_en: "Men's Straight Fit Denim Jeans",
    description_vi: "Quần denim co giãn nhẹ, phong cách cổ điển dễ phối đồ.",
    description_en: "Slightly stretchy denim jeans in a classic, easy-to-match style.",
    unit: "piece",
    base_price: NumberLong(65000000),
    sale_price: NumberLong(58000000),
    rating_avg: 4.3,
    rating_count: 62,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Size 30 / Dark Blue",
        price_delta: NumberLong(0),
        sku: "JEANS-008-30-DBL",
      },
      {
        id: new ObjectId(),
        variant_label: "Size 32 / Dark Blue",
        price_delta: NumberLong(0),
        sku: "JEANS-008-32-DBL",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/jeans-8.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 9. Electronics Product
  {
    category_id: electronicsId,
    sku: "MOUSE-009",
    name_vi: "Chuột máy tính không dây Ergonomic",
    name_en: "Ergonomic Wireless Computer Mouse",
    description_vi: "Thiết kế công phỏng học giảm mỏi cổ tay, kết nối Bluetooth & 2.4G.",
    description_en: "Ergonomic design reducing wrist strain, dual Bluetooth & 2.4G connection.",
    unit: "piece",
    base_price: NumberLong(45000000),
    sale_price: NumberLong(39000000),
    rating_avg: 4.8,
    rating_count: 405,
    is_active: true,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Color: Graphite",
        price_delta: NumberLong(0),
        sku: "MOUSE-009-GPH",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/mouse-9.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 10. Inactive Product (Useful for testing is_active filters)
  {
    category_id: electronicsId,
    sku: "KEYBOARD-010",
    name_vi: "Bàn phím cơ không dây TKL (Tạm ngưng bán)",
    name_en: "Wireless TKL Mechanical Keyboard (Inactive)",
    description_vi: "Bàn phím cơ switch Custom, đèn nền RGB nhiều chế độ.",
    description_en: "Custom switch mechanical keyboard with multi-mode RGB backlighting.",
    unit: "piece",
    base_price: NumberLong(180000000),
    sale_price: NumberLong(180000000),
    rating_avg: 0.0,
    rating_count: 0,
    is_active: false,
    variants: [
      {
        id: new ObjectId(),
        variant_label: "Switch: Red",
        price_delta: NumberLong(0),
        sku: "KEYBOARD-010-RED",
      },
    ],
    images: [
      {
        id: new ObjectId(),
        url: "https://cdn.example.com/products/keyboard-10.jpg",
        sort_order: 1,
      },
    ],
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("=== Catalog Database successfully seeded ===");
