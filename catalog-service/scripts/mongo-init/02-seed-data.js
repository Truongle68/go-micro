db = db.getSiblingDB("catalog-db");

// Clear existing collections to prevent duplicate key errors
db.categories.drop();
db.products.drop();

// Define Category ObjectIDs
const freshProduceId = new ObjectId();
const organicVegetablesId = new ObjectId();
const localFruitsId = new ObjectId();
const meatSeafoodId = new ObjectId();

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
      "https://img.lb.wbmdstatic.com/vim/live/webmd/consumer_assets/site_images/articles/health_tools/foods_for_circulation_slideshow/1800ss_getty_rf_spinach.jpg?resize=750px:*&output-quality=75",
      "https://www.gettystewart.com/wp-content/uploads/2014/01/spinach-fresh-in-bowl-sq.jpg",
    ],
    option_types: [{ name: "Trọng lượng", values: ["250g", "500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SPINACH-ORGANIC-250G",
        attributes: { "Trọng lượng": "250g" },
        price: { amount: NumberLong(22000), currency: "VND" },
        images:
          "https://img.lb.wbmdstatic.com/vim/live/webmd/consumer_assets/site_images/articles/health_tools/foods_for_circulation_slideshow/1800ss_getty_rf_spinach.jpg?resize=750px:*&output-quality=75",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SPINACH-ORGANIC-500G",
        attributes: { "Trọng lượng": "500g" },
        price: { amount: NumberLong(40000), currency: "VND" },
        images:
          "https://www.gettystewart.com/wp-content/uploads/2014/01/spinach-fresh-in-bowl-sq.jpg",
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
      "https://5.imimg.com/data5/ANDROID/Default/2022/2/GJ/ZY/SY/29350013/product-jpeg-500x500.jpg",
      "https://freshfocusfoodstuff.com/wp-content/uploads/2026/03/Dragon-Fruit-White-Box.webp",
      "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ4-u5Vmpkgd4WBUVgJ2wpC4zm306cwHylKuRQUyhnQX7AxcJZGu8WDRaAe&s=10",
      "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTc1KMf_RtRDmdGU2o6u--npwjoDGT_nZ3gxLwNmIkRPkTiFQopiFdhFtRA&s=10",
    ],
    option_types: [
      { name: "Loại", values: ["Ruột trắng", "Ruột đỏ"] },
      { name: "Đóng gói", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { Loại: "Ruột trắng", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(38000), currency: "VND" },
        images:
          "https://5.imimg.com/data5/ANDROID/Default/2022/2/GJ/ZY/SY/29350013/product-jpeg-500x500.jpg",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-3KG",
        attributes: { Loại: "Ruột trắng", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(110000), currency: "VND" },
        images:
          "https://freshfocusfoodstuff.com/wp-content/uploads/2026/03/Dragon-Fruit-White-Box.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { Loại: "Ruột đỏ", "Đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(100000), currency: "VND" },
        images:
          "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ4-u5Vmpkgd4WBUVgJ2wpC4zm306cwHylKuRQUyhnQX7AxcJZGu8WDRaAe&s=10",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { Loại: "Ruột đỏ", "Đóng gói": "Hộp quà 3kg" },
        price: { amount: NumberLong(90000), currency: "VND" },
        images:
          "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTc1KMf_RtRDmdGU2o6u--npwjoDGT_nZ3gxLwNmIkRPkTiFQopiFdhFtRA&s=10",
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
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
      "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
      "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
    ],
    option_types: [{ name: "Quy cách", values: ["Khay 300g", "Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images:
          "https://www.lottemart.vn/media/catalog/product/cache/0x0/2/0/2082590000002-3.jpg.webp",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images:
          "https://asdagroceries.scene7.com/is/image/asdagroceries/5054781380139",
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
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Product 4: Whole Ri6 Durian
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
      "https://www.foodie.com/img/gallery/what-is-durian-and-how-do-you-eat-it/intro-1744387044.jpg",
    ],
    option_types: [],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DURIAN-RI6-WHOLE",
        attributes: {},
        price: { amount: NumberLong(180000), currency: "VND" },
        images:
          "https://www.foodie.com/img/gallery/what-is-durian-and-how-do-you-eat-it/intro-1744387044.jpg",
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
    status: "inactive",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("==========================================================");
print("=== Seed Data Matched to BSON Models (Tươi Market) ===");
print("==========================================================");
