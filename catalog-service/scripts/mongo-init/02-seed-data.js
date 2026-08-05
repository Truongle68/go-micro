db = db.getSiblingDB("catalog-db");

// Clear existing collections to prevent duplicate key errors
db.categories.drop();
db.products.drop();

// Define Category ObjectIDs
const freshProduceId = new ObjectId();
const organicVegetablesId = new ObjectId();
const localFruitsId = new ObjectId();
const meatSeafoodId = new ObjectId();

// Warehouse IDs
const whSgn01 = "WH-SGN-01";
const whSgn02 = "WH-SGN-02";

// -----------------------------------------------------------------------------
// 1. SEED CATEGORIES
// -----------------------------------------------------------------------------
db.categories.insertMany([
  {
    _id: freshProduceId,
    parent_id: null,
    name: "Rau Củ Tươi",
    name_translation: {
      vi: "Rau Củ Tươi",
      en: "Fresh Produce",
    },
    slug: "rau-cu-tuoi",
    icon: "leaf",
    sort_order: 10,
    is_active: true,
    ancestors: [],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: organicVegetablesId,
    parent_id: freshProduceId,
    name: "Rau Hữu Cơ",
    name_translation: {
      vi: "Rau Hữu Cơ",
      en: "Organic Vegetables",
    },
    slug: "rau-huu-co",
    icon: "carrot",
    sort_order: 10,
    is_active: true,
    ancestors: [{ id: freshProduceId, name: "Rau Củ Tươi" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: localFruitsId,
    parent_id: freshProduceId,
    name: "Trái Cây Nội Địa",
    name_translation: {
      vi: "Trái Cây Nội Địa",
      en: "Local Fruits",
    },
    slug: "trai-cay-noi-dia",
    icon: "apple",
    sort_order: 20,
    is_active: true,
    ancestors: [{ id: freshProduceId, name: "Rau Củ Tươi" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: meatSeafoodId,
    parent_id: null,
    name: "Thịt & Hải Sản Tươi",
    name_translation: {
      vi: "Thịt & Hải Sản Tươi",
      en: "Fresh Meat & Seafood",
    },
    slug: "thit-hai-san-tuoi",
    icon: "fish",
    sort_order: 20,
    is_active: true,
    ancestors: [],
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

// Breadcrumbs matching categoryRefModel
const producePath = [{ id: freshProduceId, name: "Rau Củ Tươi" }];

const organicVegPath = [
  ...producePath,
  { id: organicVegetablesId, name: "Rau Hữu Cơ" },
];

const localFruitsPath = [
  ...producePath,
  { id: localFruitsId, name: "Trái Cây Nội Địa" },
];

const meatSeafoodPath = [{ id: meatSeafoodId, name: "Thịt & Hải Sản Tươi" }];

// -----------------------------------------------------------------------------
// 2. SEED PRODUCTS
// -----------------------------------------------------------------------------
db.products.insertMany([
  // Product 1: Da Lat Organic Spinach
  {
    _id: new ObjectId(),
    version: 1,
    slug: "rau-cai-bo-xoi-huu-co-da-lat",
    name: "Rau Cải Bó Xôi Hữu Cơ Đà Lạt",
    name_translation: {
      vi: "Rau Cải Bó Xôi Hữu Cơ Đà Lạt",
      en: "Da Lat Organic Spinach",
    },
    category_id: organicVegetablesId,
    category_path: organicVegPath,
    description:
      "Cải bó xôi (rau bina) trồng thủy canh tiêu chuẩn VietGAP tại Đà Lạt. Giàu sắt, vitamin A và C.",
    description_html:
      '<div class="product-description space-y-6 text-[#16422F]"><p class="text-base leading-relaxed font-medium"><strong>Cải bó xôi tươi (Rau bina)</strong> được thu hoạch thủ công ngay tại vườn Đà Lạt lúc sáng sớm và giao trong ngày. Lá cải xanh mướt, thân mọng nước, giàu hàm lượng Sắt, Canxi và Vitamin A—sự lựa chọn hoàn hảo cho bữa ăn gia đình và thực đơn ăn dặm của bé.</p><div class="bg-[#FAF6EC] border border-[#EBE6DA] rounded-2xl p-4 sm:p-5"><h4 class="text-sm font-extrabold text-[#1B4D3E] uppercase tracking-wide mb-3 flex items-center gap-2">🌱 Cam kết chất lượng từ Tươi Market</h4><ul class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 text-xs font-semibold text-[#2C5E4A]"><li class="flex items-center gap-2"><span class="text-emerald-600">✓</span> Thu hoạch tươi mới trong ngày</li><li class="flex items-center gap-2"><span class="text-emerald-600">✓</span> Trồng theo hướng hữu cơ tại Đà Lạt</li><li class="flex items-center gap-2"><span class="text-emerald-600">✓</span> Không thuốc trừ sâu & chất bảo quản</li><li class="flex items-center gap-2"><span class="text-emerald-600">✓</span> An toàn tuyệt đối cho bé ăn dặm</li></ul></div><div class="space-y-4"><div><h3 class="text-sm font-bold text-[#16422F] mb-1">💡 Gợi ý chế biến</h3><p class="text-xs text-slate-600 leading-relaxed">Thích hợp nấu cháo, xay sinh tố detox, xào tỏi, làm salad hoặc nấu canh thịt băm. Nên nấu nhanh trên lửa lớn để giữ trọn vị ngọt tự nhiên và dưỡng chất.</p></div><div><h3 class="text-sm font-bold text-[#16422F] mb-1">❄️ Hướng dẫn bảo quản</h3><p class="text-xs text-slate-600 leading-relaxed">Không rửa trước khi cho vào tủ lạnh. Bọc rau trong khăn giấy khô rồi cho vào túi zip đựng thực phẩm, bảo quản ở ngăn mát (4–8°C) giữ tươi từ 3–5 ngày.</p></div></div></div>',
    highlights: [
      "Đạt chứng nhận VietGAP",
      "Thu hoạch tươi mới mỗi sáng",
      "Không hóa chất bảo quản",
    ],
    tags: ["rau-sach", "vietgap", "da-lat", "huu-co"],
    images: [
      {
        url: "https://www.mahagro.com/cdn/shop/articles/spinach.jpg?v=1490962590",
        is_primary: true,
        sort_order: 1,
        alt_text: "Rau cải bó xôi tươi",
      },
      {
        url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT__Jl2V3jLSW6beYaXZNx8xELDtr4NwfDOT9t_9jps8zE_8TuiTyEErvw&s=10",
        is_primary: false,
        sort_order: 2,
        alt_text: "Đóng gói cải bó xôi 500g",
      },
    ],
    option_types: [{ name: "Trọng lượng", values: ["250g", "500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SPINACH-ORGANIC-250G",
        attributes: { "Trọng lượng": "250g" },
        price: { amount: NumberLong(22000), currency: "VND" },
        inventory: {
          total_available: 50,
          reserved: 2,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 30 },
            { warehouse_id: whSgn02, quantity: 20 },
          ],
        },
        weight_grams: 250,
        images: [
          {
            url: "https://media-cdn-v2.laodong.vn/storage/newsportal/2026/4/4/1680092/Rau.jpg",
            is_primary: true,
            sort_order: 1,
            alt_text: "Gói cải bó xôi 250g",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SPINACH-ORGANIC-500G",
        attributes: { "Trọng lượng": "500g" },
        price: { amount: NumberLong(40000), currency: "VND" },
        inventory: {
          total_available: 30,
          reserved: 0,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 15 },
            { warehouse_id: whSgn02, quantity: 15 },
          ],
        },
        weight_grams: 500,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRvQtIxAWIiKnsyXXd-3moYbqXZVOJvxSFb_KDHLiEc67K6myOf-POfTJ9A&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Gói cải bó xôi 500g",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin xuất xứ",
        items: [
          { label: "Nơi sản xuất", value: "Đà Lạt, Lâm Đồng" },
          { label: "Tiêu chuẩn", value: "VietGAP / Organic" },
        ],
      },
      {
        group: "Bảo quản",
        items: [
          { label: "Nhiệt độ", value: "4°C - 8°C trong ngăn mát tủ lạnh" },
          { label: "Hạn sử dụng", value: "3 - 5 ngày" },
        ],
      },
    ],
    rating_summary: {
      average: 4.9,
      count: 48,
      breakdown: { 1: 0, 2: 0, 3: 0, 4: 4, 5: 44 },
    },
    sales_count: NumberLong(320),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn01, region: "Ho Chi Minh City" },
      fragile: true,
      shipping_class: "cold_express",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Product 2: Binh Thuan Red Flesh Dragon Fruit
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Product 3: Fresh Norwegian Salmon Fillet
  {
    _id: new ObjectId(),
    version: 1,
    slug: "than-ca-hoi-tuoi-phile-na-uy",
    name: "Thăn Cá Hồi Tươi Phile Na Uy",
    name_translation: {
      vi: "Thăn Cá Hồi Tươi Phile Na Uy",
      en: "Fresh Norwegian Salmon Fillet",
    },
    category_id: meatSeafoodId,
    category_path: meatSeafoodPath,
    description:
      "Cá hồi Na Uy nhập khẩu hàng không về trong ngày. Thăn phile lọc sạch xương, phù hợp làm Sashimi hoặc áp chảo.",
    description_html:
      "<p>Cá hồi tươi đạt chuẩn ăn sống (Sashimi grade). Dày thịt, ngậy thơm béo ngậy giàu Omega-3.</p>",
    highlights: [
      "Nhập khẩu hàng không tươi mới mỗi ngày",
      "Đạt tiêu chuẩn ăn Sashimi trực tiếp",
      "Miễn phí cắt thái theo yêu cầu",
    ],
    tags: ["hai-san", "ca-hoi", "sashimi", "nhap-khau"],
    images: [
      {
        url: "https://cdn.tuoimarket.vn/products/salmon-1.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thăn cá hồi phile tươi",
      },
      {
        url: "https://cdn.tuoimarket.vn/products/salmon-2.jpg",
        is_primary: false,
        sort_order: 2,
        alt_text: "Cá hồi áp chảo",
      },
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        inventory: {
          total_available: 20,
          reserved: 4,
          warehouses: [{ warehouse_id: whSgn01, quantity: 20 }],
        },
        weight_grams: 300,
        images: [
          {
            url: "https://cdn.tuoimarket.vn/products/salmon-1.jpg",
            is_primary: true,
            sort_order: 1,
            alt_text: "Khay cá hồi 300g",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        inventory: {
          total_available: 15,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 15 }],
        },
        weight_grams: 500,
        images: [
          {
            url: "https://cdn.tuoimarket.vn/products/salmon-1.jpg",
            is_primary: true,
            sort_order: 1,
            alt_text: "Khay cá hồi 500g",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin hải sản",
        items: [
          { label: "Nguồn gốc", value: "Na Uy (Norwegian Atlantic Salmon)" },
          { label: "Bảo quản", value: "Giao kèm đá gel, giữ mát 0°C - 2°C" },
        ],
      },
    ],
    rating_summary: {
      average: 5.0,
      count: 112,
      breakdown: { 1: 0, 2: 0, 3: 0, 4: 0, 5: 112 },
    },
    sales_count: NumberLong(540),
    shipping: {
      is_free_shipping: true,
      ships_from: { warehouse_id: whSgn01, region: "Ho Chi Minh City" },
      fragile: true,
      shipping_class: "cold_express",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Product 4: Seasonal / Out of Stock (Sầu Riêng Ri6)
  {
    _id: new ObjectId(),
    version: 1,
    slug: "sau-rieng-ri6-viet-gap-nguyen-trai",
    name: "Sầu Riêng Ri6 VietGAP Nguyên Trái (Hết Mùa)",
    name_translation: {
      vi: "Sầu Riêng Ri6 VietGAP Nguyên Trái",
      en: "Whole Ri6 Durian (Out of Season)",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description: "Sầu riêng Ri6 cơm vàng hạt lép, béo ngậy dẻo ngọt.",
    description_html:
      "<p>Sầu riêng chín cây tự nhiên không rụng cuống. Hiện đang hết mùa thu hoạch.</p>",
    highlights: ["Cơm vàng dẻo hạt lép", "Bao đổi trả nếu bị sượng"],
    tags: ["trai-cay", "sau-rieng", "ri6", "dac-san"],
    images: [
      {
        url: "https://cdn.tuoimarket.vn/products/durian-1.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Sầu riêng Ri6",
      },
    ],
    option_types: [],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DURIAN-RI6-WHOLE",
        attributes: {},
        price: { amount: NumberLong(180000), currency: "VND" },
        inventory: {
          total_available: 0,
          reserved: 0,
          warehouses: [{ warehouse_id: whSgn02, quantity: 0 }],
        },
        weight_grams: 2500,
        images: [
          {
            url: "https://cdn.tuoimarket.vn/products/durian-1.jpg",
            is_primary: true,
            sort_order: 1,
            alt_text: "Sầu riêng nguyên trái",
          },
        ],
        is_active: false,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin",
        items: [{ label: "Xuất xứ", value: "Vĩnh Long, Việt Nam" }],
      },
    ],
    rating_summary: {
      average: 4.6,
      count: 15,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 3, 5: 11 },
    },
    sales_count: NumberLong(90),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: true,
      shipping_class: "standard",
    },
    status: "inactive",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long ngọt đậm, trái to mọng nước, giàu chất chống oxy hóa.",
    description_html:
      "<p>Thanh long tươi hái tận vườn Bình Thuận, vỏ mỏng mọng, vị ngọt tự nhiên thanh mát.</p>",
    highlights: [
      "Trái to từ 500g - 700g/trái",
      "Vị ngọt tự nhiên đậm đà",
      "Tốt cho hệ tiêu hóa & tim mạch",
    ],
    tags: ["trai-cay", "binh-thuan", "thanh-long", "giai-nhiet"],
    images: [
      {
        url: "https://product.hstatic.net/200000886243/product/dragon-fruit-pitaya-white-exoticfruitscouk-988200_7db64a23831249e3b101d64bf76d271e.jpg",
        is_primary: true,
        sort_order: 1,
        alt_text: "Thanh long",
      },
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        inventory: {
          total_available: 100,
          reserved: 5,
          warehouses: [
            { warehouse_id: whSgn01, quantity: 60 },
            { warehouse_id: whSgn02, quantity: 40 },
          ],
        },
        weight_grams: 1000,
        images: [
          {
            url: "https://nyrafoods.com/wp-content/uploads/2024/11/product-8.webp",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { "Loại": "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://cdn.jwplayer.com/v2/media/aYRMP4Wy/poster.jpg?width=720",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        inventory: {
          total_available: 35,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTes3anQ_9ob71Mc9qXCv6pf7z3vUquB3fasd0EBVIbYNRDzPdgWxvxoZw&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Túi thanh long 1kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại": "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        inventory: {
          total_available: 25,
          reserved: 1,
          warehouses: [{ warehouse_id: whSgn01, quantity: 25 }],
        },
        weight_grams: 3000,
        images: [
          {
            url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSzVL7t6ZsxDSaoizaPlRC2PhMl1jSQr1J8GDV4mSY8Jvk7es3fjvZy53M&s=10",
            is_primary: true,
            sort_order: 1,
            alt_text: "Hộp quà thanh long 3kg",
          },
        ],
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thông tin sản phẩm",
        items: [
          { label: "Xuất xứ", value: "Bình Thuận, Việt Nam" },
          { label: "Độ đường (Brix)", value: "> 14%" },
        ],
      },
    ],
    rating_summary: {
      average: 4.7,
      count: 29,
      breakdown: { 1: 0, 2: 0, 3: 1, 4: 6, 5: 22 },
    },
    sales_count: NumberLong(185),
    shipping: {
      is_free_shipping: false,
      ships_from: { warehouse_id: whSgn02, region: "Ho Chi Minh City" },
      fragile: false,
      shipping_class: "standard",
    },
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("==========================================================");
print("=== Seed Data Matched to BSON Models (Tươi Market) ===");
print("==========================================================");
