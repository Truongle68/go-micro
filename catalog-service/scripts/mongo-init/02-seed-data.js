db = db.getSiblingDB("catalog-db");

// Clear existing collections to prevent duplicate key errors
db.categories.drop();
db.products.drop();

// -----------------------------------------------------------------------------
// 1. DEFINE CATEGORY IDS & SEED HIERARCHY
// -----------------------------------------------------------------------------
const rootProduceId = new ObjectId();
const rootMeatFishId = new ObjectId();
const rootDairyId = new ObjectId();

const organicVegId = new ObjectId();
const rootVegId = new ObjectId();
const localFruitsId = new ObjectId();
const importedFruitsId = new ObjectId();
const freshMeatId = new ObjectId();
const seafoodId = new ObjectId();

db.categories.insertMany([
  // Root Categories
  {
    _id: rootProduceId,
    parent_id: null,
    name: "Rau Củ Quả",
    name_translation: { vi: "Rau Củ Quả", en: "Fresh Produce" },
    slug: "rau-cu-qua",
    icon: "leaf",
    sort_order: 10,
    is_active: true,
    ancestors: [],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: rootMeatFishId,
    parent_id: null,
    name: "Thịt & Hải Sản",
    name_translation: { vi: "Thịt & Hải Sản", en: "Meat & Seafood" },
    slug: "thit-hai-san",
    icon: "fish",
    sort_order: 20,
    is_active: true,
    ancestors: [],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: rootDairyId,
    parent_id: null,
    name: "Trứng & Sữa Tươi",
    name_translation: { vi: "Trứng & Sữa Tươi", en: "Dairy & Eggs" },
    slug: "trung-sua-tuoi",
    icon: "egg",
    sort_order: 30,
    is_active: true,
    ancestors: [],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Sub-categories: Produce
  {
    _id: organicVegId,
    parent_id: rootProduceId,
    name: "Rau Lá Hữu Cơ",
    name_translation: { vi: "Rau Lá Hữu Cơ", en: "Organic Leafy Greens" },
    slug: "rau-la-huu-co",
    icon: "sprout",
    sort_order: 10,
    is_active: true,
    ancestors: [{ id: rootProduceId, name: "Rau Củ Quả" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: rootVegId,
    parent_id: rootProduceId,
    name: "Củ & Quả Nấu Ăn",
    name_translation: { vi: "Củ & Quả Nấu Ăn", en: "Root Vegetables & Gourds" },
    slug: "cu-qua-nau-an",
    icon: "carrot",
    sort_order: 20,
    is_active: true,
    ancestors: [{ id: rootProduceId, name: "Rau Củ Quả" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: localFruitsId,
    parent_id: rootProduceId,
    name: "Trái Cây Nội Địa",
    name_translation: { vi: "Trái Cây Nội Địa", en: "Domestic Fruits" },
    slug: "trai-cay-noi-dia",
    icon: "apple",
    sort_order: 30,
    is_active: true,
    ancestors: [{ id: rootProduceId, name: "Rau Củ Quả" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: importedFruitsId,
    parent_id: rootProduceId,
    name: "Trái Cây Nhập Khẩu",
    name_translation: { vi: "Trái Cây Nhập Khẩu", en: "Imported Fruits" },
    slug: "trai-cay-nhap-khau",
    icon: "globe",
    sort_order: 40,
    is_active: true,
    ancestors: [{ id: rootProduceId, name: "Rau Củ Quả" }],
    created_at: new Date(),
    updated_at: new Date(),
  },

  // Sub-categories: Meat & Seafood
  {
    _id: freshMeatId,
    parent_id: rootMeatFishId,
    name: "Thịt Tươi Sạch",
    name_translation: { vi: "Thịt Tươi Sạch", en: "Fresh Meat" },
    slug: "thit-tuoi-sach",
    icon: "drumstick",
    sort_order: 10,
    is_active: true,
    ancestors: [{ id: rootMeatFishId, name: "Thịt & Hải Sản" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: seafoodId,
    parent_id: rootMeatFishId,
    name: "Hải Sản Cao Cấp",
    name_translation: { vi: "Hải Sản Cao Cấp", en: "Seafood" },
    slug: "hai-san-cao-cap",
    icon: "anchor",
    sort_order: 20,
    is_active: true,
    ancestors: [{ id: rootMeatFishId, name: "Thịt & Hải Sản" }],
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

// Breadcrumbs matching categoryRefModel
const organicVegPath = [
  { id: rootProduceId, name: "Rau Củ Quả" },
  { id: organicVegId, name: "Rau Lá Hữu Cơ" },
];
const rootVegPath = [
  { id: rootProduceId, name: "Rau Củ Quả" },
  { id: rootVegId, name: "Củ & Quả Nấu Ăn" },
];
const localFruitsPath = [
  { id: rootProduceId, name: "Rau Củ Quả" },
  { id: localFruitsId, name: "Trái Cây Nội Địa" },
];
const importedFruitsPath = [
  { id: rootProduceId, name: "Rau Củ Quả" },
  { id: importedFruitsId, name: "Trái Cây Nhập Khẩu" },
];
const meatPath = [
  { id: rootMeatFishId, name: "Thịt & Hải Sản" },
  { id: freshMeatId, name: "Thịt Tươi Sạch" },
];
const seafoodPath = [
  { id: rootMeatFishId, name: "Thịt & Hải Sản" },
  { id: seafoodId, name: "Hải Sản Cao Cấp" },
];
const dairyPath = [{ id: rootDairyId, name: "Trứng & Sữa Tươi" }];

// -----------------------------------------------------------------------------
// 2. SEED 20 DIVERSE PRODUCTS
// -----------------------------------------------------------------------------
db.products.insertMany([
  // 1. Spinach
  {
    _id: new ObjectId(),
    version: 1,
    slug: "rau-cai-bo-xoi-huu-co-da-lat",
    name: "Rau Cải Bó Xôi Hữu Cơ Đà Lạt",
    name_translation: {
      vi: "Rau Cải Bó Xôi Hữu Cơ Đà Lạt",
      en: "Da Lat Organic Baby Spinach",
    },
    category_id: organicVegId,
    category_path: organicVegPath,
    description:
      "Cải bó xôi (rau bina) trồng chuẩn VietGAP tại Đà Lạt. Giàu sắt và khoáng chất.",
    description_html:
      "<p>Cải bó xôi hái tươi tại vườn, lá xanh non, an toàn chuẩn hữu cơ.</p>",
    highlights: ["Chuẩn VietGAP", "Không chất bảo quản", "Thu hoạch sớm"],
    tags: ["rau-sach", "vietgap", "da-lat", "huu-co"],
    images: ["https://images.unsplash.com/photo-1576045057995-568f588f82fb"],
    option_types: [{ name: "Trọng lượng", values: ["250g", "500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SPINACH-250G",
        attributes: { "Trọng lượng": "250g" },
        price: { amount: NumberLong(22000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1576045057995-568f588f82fb",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SPINACH-500G",
        attributes: { "Trọng lượng": "500g" },
        price: { amount: NumberLong(40000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1576045057995-568f588f82fb",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Xuất xứ",
        items: [{ label: "Địa phương", value: "Đà Lạt, Lâm Đồng" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 2. Kale
  {
    _id: new ObjectId(),
    version: 1,
    slug: "cai-kale-xoan-thuy-canh",
    name: "Cải Kale Xoăn Thủy Canh",
    name_translation: {
      vi: "Cải Kale Xoăn Thủy Canh",
      en: "Curly Hydroponic Kale",
    },
    category_id: organicVegId,
    category_path: organicVegPath,
    description:
      "Cải kale giàu canxi và vitamin K, chuyên dùng ép nước detox hoặc làm salad giòn.",
    description_html:
      "<p>Kale xoăn hữu cơ giàu dinh dưỡng, cuống non giòn sần sật.</p>",
    highlights: ["Siêu thực phẩm (Superfood)", "Chuyên detox và ăn kiêng"],
    tags: ["kale", "detox", "an-kieng", "thuy-canh"],
    images: ["https://images.unsplash.com/photo-1524179091875-bf99a9a6af57"],
    option_types: [{ name: "Trọng lượng", values: ["300g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "KALE-300G",
        attributes: { "Trọng lượng": "300g" },
        price: { amount: NumberLong(32000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1524179091875-bf99a9a6af57",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Canh tác",
        items: [{ label: "Phương pháp", value: "Thủy canh tuần hoàn" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 3. Bok Choy
  {
    _id: new ObjectId(),
    version: 1,
    slug: "cai-thi-ha-noi-baby",
    name: "Cải Thìa Baby Thượng Hạng",
    name_translation: {
      vi: "Cải Thìa Baby Thượng Hạng",
      en: "Baby Bok Choy",
    },
    category_id: organicVegId,
    category_path: organicVegPath,
    description:
      "Cải thìa non cuống dày mọng nước, vị ngọt thanh mát khi xào tỏi hoặc nấu súp.",
    description_html:
      "<p>Cải thìa giống baby giòn ngọt tự nhiên, không xơ ráp.</p>",
    highlights: ["Không dư lượng thuốc trừ sâu", "Thu hoạch dạng cây non"],
    tags: ["cai-thia", "rau-sach", "vietgap"],
    images: ["https://images.unsplash.com/photo-1587049352846-4a222e784d38"],
    option_types: [{ name: "Trọng lượng", values: ["500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "BOKCHOY-500G",
        attributes: { "Trọng lượng": "500g" },
        price: { amount: NumberLong(18000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1587049352846-4a222e784d38",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Bảo quản",
        items: [{ label: "Nhiệt độ", value: "3°C - 6°C" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 4. Carrots
  {
    _id: new ObjectId(),
    version: 1,
    slug: "ca-rot-baby-da-lat",
    name: "Cà Rốt Baby Đà Lạt",
    name_translation: {
      vi: "Cà Rốt Baby Đà Lạt",
      en: "Da Lat Mini Carrots",
    },
    category_id: rootVegId,
    category_path: rootVegPath,
    description:
      "Cà rốt giống tí hon nguyên cuống xanh, vị đậm đà không xơ, thích hợp nướng mật ong.",
    description_html:
      "<p>Cà rốt baby giòn ngọt tươi ngon mới nhổ từ đất đỏ bazan.</p>",
    highlights: ["Nguyên củ có cuống lá tươi", "Đậm đà vị đất lạnh tự nhiên"],
    tags: ["ca-rot", "cu-qua", "da-lat"],
    images: ["https://images.unsplash.com/photo-1598170845058-32b9d6a5da37"],
    option_types: [{ name: "Quy cách", values: ["Bó 500g", "Túi 1kg"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "CARROT-MINI-500G",
        attributes: { "Quy cách": "Bó 500g" },
        price: { amount: NumberLong(28000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1598170845058-32b9d6a5da37",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "CARROT-MINI-1KG",
        attributes: { "Quy cách": "Túi 1kg" },
        price: { amount: NumberLong(52000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1598170845058-32b9d6a5da37",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Xuất xứ",
        items: [{ label: "Địa chỉ vườn", value: "Đơn Dương, Lâm Đồng" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 5. Japanese Sweet Potato
  {
    _id: new ObjectId(),
    version: 1,
    slug: "khoai-lang-mat-da-lat",
    name: "Khoai Lang Mật Đà Lạt Thượng Hạng",
    name_translation: {
      vi: "Khoai Lang Mật Đà Lạt Thượng Hạng",
      en: "Sweet Honey Sweet Potato",
    },
    category_id: rootVegId,
    category_path: rootVegPath,
    description:
      "Khoai lang mật dẻo thơm, nướng tươm mật ngọt đậm đà, củ thuôn đều dễ chế biến.",
    description_html:
      "<p>Khoai đã được ủ xuống mật tự nhiên để có độ ngọt cao nhất khi nướng hoặc hấp.</p>",
    highlights: ["Ủ mật tự nhiên", "Cơm vàng óng mướt"],
    tags: ["khoai-lang", "khoai-mat", "nong-san-viet"],
    images: ["https://images.unsplash.com/photo-1596097635121-14b63b7a0c19"],
    option_types: [{ name: "Trọng lượng", values: ["1kg", "2kg"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SWEETPOTATO-1KG",
        attributes: { "Trọng lượng": "1kg" },
        price: { amount: NumberLong(42000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1596097635121-14b63b7a0c19",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SWEETPOTATO-2KG",
        attributes: { "Trọng lượng": "2kg" },
        price: { amount: NumberLong(80000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1596097635121-14b63b7a0c19",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Chất lượng",
        items: [{ label: "Tỉ lệ ngọt", value: "Cao (ủ mật > 10 ngày)" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 6. Cherry Tomatoes
  {
    _id: new ObjectId(),
    version: 1,
    slug: "ca-chua-bi-soc-nau-ngot",
    name: "Cà Chua Bi Sô-cô-la Ngọt",
    name_translation: {
      vi: "Cà Chua Bi Sô-cô-la Ngọt",
      en: "Sweet Chocolate Cherry Tomatoes",
    },
    category_id: rootVegId,
    category_path: rootVegPath,
    description:
      "Cà chua bi nâu sô-cô-la mọng quả, vị ngọt hậu thanh nhẹ, giàu chất chống oxy hóa Lycopene.",
    description_html:
      "<p>Loại cà chua bi cao cấp thích hợp ăn sống như trái cây hoặc trộn salad.</p>",
    highlights: ["Độ brix > 9", "Vỏ bóng, thịt dày"],
    tags: ["ca-chua", "salad", "trai-cay-an-song"],
    images: ["https://images.unsplash.com/photo-1592924357228-91a4daadcfea"],
    option_types: [{ name: "Hộp", values: ["500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "TOMATO-CHOC-500G",
        attributes: { Hộp: "500g" },
        price: { amount: NumberLong(39000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1592924357228-91a4daadcfea",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Canh tác",
        items: [{ label: "Mô hình", value: "Nhà màng công nghệ cao" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 7. Dragon Fruit
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thanh-long-binh-thuan",
    name: "Thanh Long Bình Thuận Ruột Đỏ / Trắng",
    name_translation: {
      vi: "Thanh Long Bình Thuận",
      en: "Binh Thuan Dragon Fruit",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Thanh long tươi mọng nước, vỏ mỏng, tai xanh giòn, vị ngọt dịu nhẹ.",
    description_html:
      "<p>Đặc sản vùng đất nắng gió Bình Thuận, chuẩn xuất khẩu.</p>",
    highlights: ["Tai xanh tươi, vỏ mỏng", "Thịt quả mọng nước"],
    tags: ["thanh-long", "trai-cay", "binh-thuan"],
    images: ["https://images.unsplash.com/photo-1527325678964-54921661f888"],
    option_types: [
      { name: "Loại ruột", values: ["Ruột đỏ", "Ruột trắng"] },
      { name: "Quy cách", values: ["Túi 1kg", "Hộp quà 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-1KG",
        attributes: { "Loại ruột": "Ruột đỏ", "Quy cách": "Túi 1kg" },
        price: { amount: NumberLong(45000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1527325678964-54921661f888",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-WHITE-1KG",
        attributes: { "Loại ruột": "Ruột trắng", "Quy cách": "Túi 1kg" },
        price: { amount: NumberLong(32000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1527325678964-54921661f888",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "DRAGON-RED-3KG",
        attributes: { "Loại ruột": "Ruột đỏ", "Quy cách": "Hộp quà 3kg" },
        price: { amount: NumberLong(130000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1527325678964-54921661f888",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Xuất xứ",
        items: [{ label: "Nguồn gốc", value: "Hàm Thuận Nam, Bình Thuận" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 8. Ri6 Durian
  {
    _id: new ObjectId(),
    version: 1,
    slug: "sau-rieng-ri6-chin-cay-vinh-long",
    name: "Sầu Riêng Ri6 Chín Cây Tách Múi",
    name_translation: {
      vi: "Sầu Riêng Ri6 Chín Cây Tách Múi",
      en: "Ri6 Durian Clean Pulp Tray",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Cơm sầu riêng Ri6 rụng tự nhiên, múi vàng óng, dẻo béo đậm vị không xơ sượng.",
    description_html:
      "<p>Được khui trực tiếp từ trái sầu chín rụng tự nhiên, đóng khay hút màng bảo vệ.</p>",
    highlights: ["Bao 1 đổi 1 nếu sượng", "Cơm dẻo khô ráo"],
    tags: ["sau-rieng", "ri6", "mien-tay"],
    images: ["https://images.unsplash.com/photo-1587393855524-087f83d95bc9"],
    option_types: [{ name: "Khối lượng múi", values: ["Khay 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "DURIAN-RI6-500G",
        attributes: { "Khối lượng múi": "Khay 500g" },
        price: { amount: NumberLong(165000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1587393855524-087f83d95bc9",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Bảo hành",
        items: [{ label: "Cam kết", value: "Hư sượng đền bù 100%" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 9. Hoa Loc Mango
  {
    _id: new ObjectId(),
    version: 1,
    slug: "xoai-cat-hoa-loc-tien-giang",
    name: "Xoài Cát Hòa Lộc Tiền Giang",
    name_translation: {
      vi: "Xoài Cát Hòa Lộc Tiền Giang",
      en: "Tien Giang Hoa Loc Mangoes",
    },
    category_id: localFruitsId,
    category_path: localFruitsPath,
    description:
      "Vua của các loài xoài Việt Nam. Vỏ mỏng vàng ruộm, thịt thơm nức mũi và không chút xơ.",
    description_html:
      "<p>Xoài cát tuyển lựa từng quả từ nhà vườn Cái Bè - Tiền Giang.</p>",
    highlights: ["Độ ngọt đậm brix > 18", "Thịt dày hạt mỏng dính"],
    tags: ["xoai-cat", "hoa-loc", "trai-cay-mien-tay"],
    images: ["https://images.unsplash.com/photo-1553279768-865429fa0078"],
    option_types: [
      { name: "Kích cỡ quả", values: ["Loại 1 (450-600g/quả)", "Thùng 3kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "MANGO-HL-L1",
        attributes: { "Kích cỡ quả": "Loại 1 (450-600g/quả)" },
        price: { amount: NumberLong(85000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1553279768-865429fa0078",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "MANGO-HL-BOX3K",
        attributes: { "Kích cỡ quả": "Thùng 3kg" },
        price: { amount: NumberLong(245000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1553279768-865429fa0078",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Địa lý",
        items: [{ label: "Vùng trồng", value: "Cái Bè, Tiền Giang" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 10. Envy Apples (Imported)
  {
    _id: new ObjectId(),
    version: 1,
    slug: "tao-envy-new-zealand-size-70",
    name: "Táo Envy New Zealand Size Lớn",
    name_translation: {
      vi: "Táo Envy New Zealand Size Lớn",
      en: "New Zealand Envy Apples",
    },
    category_id: importedFruitsId,
    category_path: importedFruitsPath,
    description:
      "Táo Envy vỏ đỏ ruby rực rỡ, thịt giòn đanh, mọng nước và có mùi thơm thoảng mật ong đặc trưng.",
    description_html:
      "<p>Táo nhập khẩu tươi bằng đường hàng không, cuống còn xanh tươi.</p>",
    highlights: ["Thịt giòn, lâu ngả màu thâm", "Nhập khẩu nguyên kiện"],
    tags: ["tao-envy", "nhap-khau", "new-zealand"],
    images: ["https://images.unsplash.com/photo-1560806887-1e4cd0b6cbd6"],
    option_types: [
      { name: "Quy cách đóng gói", values: ["Khay 3 quả", "Túi 1kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "APPLE-ENVY-3PCS",
        attributes: { "Quy cách đóng gói": "Khay 3 quả" },
        price: { amount: NumberLong(95000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1560806887-1e4cd0b6cbd6",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "APPLE-ENVY-1KG",
        attributes: { "Quy cách đóng gói": "Túi 1kg" },
        price: { amount: NumberLong(165000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1560806887-1e4cd0b6cbd6",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Xuất nhập khẩu",
        items: [{ label: "Xuất xứ", value: "New Zealand" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 11. Autumn Crisp Grapes (Imported)
  {
    _id: new ObjectId(),
    version: 1,
    slug: "nho-xanh-khong-hat-autumn-crisp-uc",
    name: "Nho Xanh Không Hạt Autumn Crisp Úc",
    name_translation: {
      vi: "Nho Xanh Không Hạt Autumn Crisp Úc",
      en: "Australian Green Seedless Grapes",
    },
    category_id: importedFruitsId,
    category_path: importedFruitsPath,
    description:
      "Quả nho căng tròn, vị giòn tan rôm rốp trong miệng, ngọt đậm đà không hề gắt.",
    description_html:
      "<p>Nho xanh Autumn Crisp loại hảo hạng được tuyển chọn từ các trang trại Nam Úc.</p>",
    highlights: ["Không hạt, vỏ mỏng không chát", "Giòn rụm từng trái"],
    tags: ["nho-xanh", "nho-uc", "trai-cay-nhap"],
    images: ["https://images.unsplash.com/photo-1537640538966-79f369143f8f"],
    option_types: [{ name: "Trọng lượng", values: ["Hộp 500g", "Hộp 1kg"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "GRAPE-AC-500G",
        attributes: { "Trọng lượng": "Hộp 500g" },
        price: { amount: NumberLong(125000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1537640538966-79f369143f8f",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "GRAPE-AC-1KG",
        attributes: { "Trọng lượng": "Hộp 1kg" },
        price: { amount: NumberLong(240000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1537640538966-79f369143f8f",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Bảo quản",
        items: [{ label: "Nhiệt độ tối ưu", value: "-1°C đến 2°C" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 12. Pork Belly
  {
    _id: new ObjectId(),
    version: 1,
    slug: "thit-ba-roi-heo-sinh-hoc",
    name: "Thịt Ba Rọi Heo Quế Sinh Học",
    name_translation: {
      vi: "Thịt Ba Rọi Heo Quế Sinh Học",
      en: "Clean Pork Belly Slices",
    },
    category_id: freshMeatId,
    category_path: meatPath,
    description:
      "Thịt heo nuôi thảo mộc, nạc mỡ xen kẽ 5 lớp hoàn hảo, không chất tăng trọng.",
    description_html:
      "<p>Thịt heo tươi nóng mổ trong ngày, mỡ thơm giòn không ngấy, bì mềm dẻo.</p>",
    highlights: ["Không kháng sinh tồn dư", "Thịt mềm thơm tự nhiên"],
    tags: ["thit-heo", "ba-roi", "thit-sach"],
    images: ["https://images.unsplash.com/photo-1607623814075-e51df1bdc82f"],
    option_types: [
      { name: "Quy cách cắt", values: ["Nguyên tảng", "Thái lát mỏng"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "PORK-BELLY-WHOLE-500G",
        attributes: { "Quy cách cắt": "Nguyên tảng" },
        price: { amount: NumberLong(95000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1607623814075-e51df1bdc82f",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "PORK-BELLY-SLICED-500G",
        attributes: { "Quy cách cắt": "Thái lát mỏng" },
        price: { amount: NumberLong(98000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1607623814075-e51df1bdc82f",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "An toàn",
        items: [
          { label: "Kiểm nghiệm", value: "Không chất tạo nạc Clenbuterol" },
        ],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 13. Free-range Chicken
  {
    _id: new ObjectId(),
    version: 1,
    slug: "ga-ta-tha-vuon-lam-sach",
    name: "Gà Ta Thả Vườn Đồi Làm Sạch",
    name_translation: {
      vi: "Gà Ta Thả Vườn Đồi Làm Sạch",
      en: "Free-range Whole Chicken",
    },
    category_id: freshMeatId,
    category_path: meatPath,
    description:
      "Gà ta nuôi đồi tự nhiên, thịt săn chắc ngọt đậm, da vàng giòn sần sật, mỡ vừa phải.",
    description_html:
      "<p>Gà được làm sạch hút chân không kèm lòng mề và muối tiêu chanh ớt.</p>",
    highlights: ["Thịt chắc ngọt thơm", "Nuôi thức ăn tự nhiên ngô thóc"],
    tags: ["ga-ta", "thit-ga", "dac-san-doi"],
    images: ["https://images.unsplash.com/photo-1587593810167-a84920ea0781"],
    option_types: [
      { name: "Quy cách", values: ["Nguyên con (~1.3kg - 1.5kg)"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "CHICKEN-WHOLE-1400",
        attributes: { "Quy cách": "Nguyên con (~1.3kg - 1.5kg)" },
        price: { amount: NumberLong(185000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1587593810167-a84920ea0781",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Trang trại",
        items: [{ label: "Khu vực", value: "Yên Thế, Bắc Giang" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 14. Wagyu Beef Ribeye
  {
    _id: new ObjectId(),
    version: 1,
    slug: "dau-that-lung-bo-wagyu-uc-mb4-5",
    name: "Đầu Thắt Lưng Bò Wagyu Úc Ribeye MB4-5",
    name_translation: {
      vi: "Đầu Thắt Lưng Bò Wagyu Úc Ribeye MB4-5",
      en: "Australian Wagyu Ribeye MB4-5",
    },
    category_id: freshMeatId,
    category_path: meatPath,
    description:
      "Vân mỡ cẩm thạch hoàn mỹ, vị ngậy thơm như bơ khi áp chảo, mềm tan nơi đầu lưỡi.",
    description_html:
      "<p>Bò Wagyu ăn ngũ cốc 300 ngày, tiêu chuẩn vân mỡ Marble Score MB4/5.</p>",
    highlights: ["Vân mỡ cẩm thạch đồng đều", "Thích hợp làm Steak hảo hạng"],
    tags: ["bo-uc", "wagyu", "steak", "thit-bo"],
    images: ["https://images.unsplash.com/photo-1558030006-450675393462"],
    option_types: [
      {
        name: "Độ dày cắt",
        values: ["Cắt Steak 2cm (250g)", "Cắt Steak 3cm (350g)"],
      },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "WAGYU-RIBEYE-250G",
        attributes: { "Độ dày cắt": "Cắt Steak 2cm (250g)" },
        price: { amount: NumberLong(340000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1558030006-450675393462",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "WAGYU-RIBEYE-350G",
        attributes: { "Độ dày cắt": "Cắt Steak 3cm (350g)" },
        price: { amount: NumberLong(470000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1558030006-450675393462",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Tiêu chuẩn",
        items: [{ label: "Chỉ số vân mỡ", value: "MB4 - MB5" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 15. Norwegian Salmon
  {
    _id: new ObjectId(),
    version: 1,
    slug: "than-ca-hoi-tuoi-phile-na-uy",
    name: "Thăn Cá Hồi Tươi Phile Na Uy (Sashimi Grade)",
    name_translation: {
      vi: "Thăn Cá Hồi Tươi Phile Na Uy",
      en: "Fresh Norwegian Salmon Fillet",
    },
    category_id: seafoodId,
    category_path: seafoodPath,
    description:
      "Cá hồi Na Uy đường bay lọc xương tỉ mỉ, mọng mỡ, thơm béo ngậy chuẩn sashimi.",
    description_html:
      "<p>Phile cá hồi Nauy đánh bắt tươi trong tuần, chuẩn ăn sống trực tiếp.</p>",
    highlights: ["Sashimi Grade an toàn", "Bảo quản lạnh không đông đá"],
    tags: ["ca-hoi", "sashimi", "na-uy", "hai-san"],
    images: ["https://images.unsplash.com/photo-1499125562588-29fb8a56b5d5"],
    option_types: [
      { name: "Quy cách đóng gói", values: ["Khay 300g", "Khay 500g"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-300G",
        attributes: { "Quy cách đóng gói": "Khay 300g" },
        price: { amount: NumberLong(195000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1499125562588-29fb8a56b5d5",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SALMON-NOR-500G",
        attributes: { "Quy cách đóng gói": "Khay 500g" },
        price: { amount: NumberLong(320000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1499125562588-29fb8a56b5d5",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Hải sản",
        items: [{ label: "Nguồn gốc", value: "Salmar / Lerøy (Na Uy)" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 16. Black Tiger Shrimp
  {
    _id: new ObjectId(),
    version: 1,
    slug: "tom-su-bien-song-ca-mau",
    name: "Tôm Sú Biển Thiên Nhiên Cà Mau",
    name_translation: {
      vi: "Tôm Sú Biển Thiên Nhiên Cà Mau",
      en: "Ca Mau Wild Black Tiger Prawns",
    },
    category_id: seafoodId,
    category_path: seafoodPath,
    description:
      "Tôm sú rừng ngập mặn vỏ mỏng giòn, thịt dai chắc ngọt nước tự nhiên.",
    description_html:
      "<p>Đánh bắt thủ công tại các vuông tôm sinh thái ngập mặn Cà Mau.</p>",
    highlights: [
      "Vỏ mỏng thịt dai ngọt đậm",
      "Cấp đông sâu IQF giữ trọn dinh dưỡng",
    ],
    tags: ["tom-su", "ca-mau", "hai-san-tuoi"],
    images: ["https://images.unsplash.com/photo-1565680018434-b513d5e5fd47"],
    option_types: [
      { name: "Kích cỡ", values: ["Size 20-25 con/kg", "Size 10-15 con/kg"] },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "SHRIMP-SZ20-500G",
        attributes: { "Kích cỡ": "Size 20-25 con/kg" },
        price: { amount: NumberLong(145000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1565680018434-b513d5e5fd47",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "SHRIMP-SZ10-500G",
        attributes: { "Kích cỡ": "Size 10-15 con/kg" },
        price: { amount: NumberLong(210000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1565680018434-b513d5e5fd47",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Môi trường",
        items: [{ label: "Hình thức", value: "Sinh thái rừng đước" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 17. Live Mud Crab
  {
    _id: new ObjectId(),
    version: 1,
    slug: "cua-bien-ca-mau-song-day-khong-trong-luong",
    name: "Cua Cà Mau Dây Không Trọng Lượng",
    name_translation: {
      vi: "Cua Cà Mau Dây Không Trọng Lượng",
      en: "Ca Mau Live Mud Crab",
    },
    category_id: seafoodId,
    category_path: seafoodPath,
    description:
      "Cua biển Cà Mau chắc nịch, dây trói mỏng không ngậm nước, bao ăn 1 đổi 1 tận nhà.",
    description_html:
      "<p>Cua đực thịt ngọt chắc hoặc cua cái đầy gạch son béo ngậy.</p>",
    highlights: ["Dây vải cực nhỏ không ăn gian cân", "Bao ốp 100%"],
    tags: ["cua-ca-mau", "cua-thit", "cua-gach"],
    images: ["https://images.unsplash.com/photo-1559742811-82286364ceaf"],
    option_types: [
      {
        name: "Loại cua",
        values: ["Cua Y (Thịt) 400g-500g", "Cua Gạch 400g-500g"],
      },
    ],
    variants: [
      {
        _id: new ObjectId(),
        sku: "CRAB-MEAT-450G",
        attributes: { "Loại cua": "Cua Y (Thịt) 400g-500g" },
        price: { amount: NumberLong(220000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1559742811-82286364ceaf",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "CRAB-ROE-450G",
        attributes: { "Loại cua": "Cua Gạch 400g-500g" },
        price: { amount: NumberLong(310000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1559742811-82286364ceaf",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Cam kết",
        items: [{ label: "Chất lượng thịt", value: "Độ chắc thịt trên 90%" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 18. Organic Chicken Eggs
  {
    _id: new ObjectId(),
    version: 1,
    slug: "trung-ga-ta-an-thao-moc",
    name: "Trứng Gà Thảo Mộc Omega-3",
    name_translation: {
      vi: "Trứng Gà Thảo Mộc Omega-3",
      en: "Herbal-Fed Omega-3 Chicken Eggs",
    },
    category_id: rootDairyId,
    category_path: dairyPath,
    description:
      "Trứng gà lòng đỏ cam sậm đậm đà, không tanh, giàu Omega-3 và vitamin E.",
    description_html:
      "<p>Đàn gà đẻ được cho ăn thức ăn lên men từ thảo dược tự nhiên và dầu cá hồi.</p>",
    highlights: ["Lòng đỏ dày béo ngậy", "Chuẩn an toàn ăn lòng đào"],
    tags: ["trung-ga", "omega3", "organic"],
    images: ["https://images.unsplash.com/photo-1506976785307-8732e854ad03"],
    option_types: [{ name: "Quy cách vỉ", values: ["Vỉ 6 quả", "Hộp 10 quả"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "EGGS-OMEGA-6PCS",
        attributes: { "Quy cách vỉ": "Vỉ 6 quả" },
        price: { amount: NumberLong(31000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1506976785307-8732e854ad03",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "EGGS-OMEGA-10PCS",
        attributes: { "Quy cách vỉ": "Hộp 10 quả" },
        price: { amount: NumberLong(49000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1506976785307-8732e854ad03",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Dinh dưỡng",
        items: [{ label: "Omega-3", value: "> 120mg / quả" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 19. Pasteurized Fresh Milk
  {
    _id: new ObjectId(),
    version: 1,
    slug: "sua-tuoi-thanh-trung-nguyen-chat-da-lat",
    name: "Sữa Tươi Thanh Trùng Nguyên Chất Đà Lạt",
    name_translation: {
      vi: "Sữa Tươi Thanh Trùng Nguyên Chất Đà Lạt",
      en: "Da Lat Pure Pasteurized Whole Milk",
    },
    category_id: rootDairyId,
    category_path: dairyPath,
    description:
      "Sữa bò nguyên chất thanh trùng ở nhiệt độ thấp, giữ nguyên vẹn lớp váng sữa béo thơm.",
    description_html:
      "<p>100% sữa tươi từ nông trại cao nguyên mát lạnh, hạn dùng 10 ngày tươi mới.</p>",
    highlights: ["Không hoàn nguyên sữa bột", "Giữ trọn kháng thể tự nhiên"],
    tags: ["sua-tuoi", "thanh-trung", "da-lat"],
    images: ["https://images.unsplash.com/photo-1550583724-b2692b85b150"],
    option_types: [{ name: "Loại đường", values: ["Không đường", "Có đường"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "MILK-UNSWEET-900ML",
        attributes: { "Loại đường": "Không đường" },
        price: { amount: NumberLong(38000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1550583724-b2692b85b150",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "MILK-SWEET-900ML",
        attributes: { "Loại đường": "Có đường" },
        price: { amount: NumberLong(38000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1550583724-b2692b85b150",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Bảo quản",
        items: [{ label: "Hạn sử dụng", value: "10 ngày ở 2°C - 4°C" }],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },

  // 20. Greek Yogurt
  {
    _id: new ObjectId(),
    version: 1,
    slug: "sua-chua-hy-lap-thu-cong-greek-yogurt",
    name: "Sữa Chua Hy Lạp Thủ Công Nguyên Kem",
    name_translation: {
      vi: "Sữa Chua Hy Lạp Thủ Công Nguyên Kem",
      en: "Artisanal Full Cream Greek Yogurt",
    },
    category_id: rootDairyId,
    category_path: dairyPath,
    description:
      "Sữa chua lên men tự nhiên được lọc tách nước whey 3 lần, sánh đặc như kem bơ, giàu protein.",
    description_html:
      "<p>Thực phẩm lý tưởng cho người tập gym, ăn eat-clean và cải thiện hệ vi sinh đường ruột.</p>",
    highlights: [
      "Hàm lượng protein cao gấp 2.5 lần",
      "Đặc dẻo tự nhiên không thạch gelatin",
    ],
    tags: ["greek-yogurt", "sua-chua", "eat-clean"],
    images: ["https://images.unsplash.com/photo-1488477181946-6428a0291777"],
    option_types: [{ name: "Khối lượng", values: ["Hũ 200g", "Hũ 500g"] }],
    variants: [
      {
        _id: new ObjectId(),
        sku: "GYOGURT-200G",
        attributes: { "Khối lượng": "Hũ 200g" },
        price: { amount: NumberLong(42000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1488477181946-6428a0291777",
        is_active: true,
        created_at: new Date(),
      },
      {
        _id: new ObjectId(),
        sku: "GYOGURT-500G",
        attributes: { "Khối lượng": "Hũ 500g" },
        price: { amount: NumberLong(95000), currency: "VND" },
        images: "https://images.unsplash.com/photo-1488477181946-6428a0291777",
        is_active: true,
        created_at: new Date(),
      },
    ],
    specifications: [
      {
        group: "Thành phần",
        items: [
          {
            label: "Men sống",
            value: "Lactobacillus bulgaricus, Streptococcus thermophilus",
          },
        ],
      },
    ],
    status: "active",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("==========================================================");
print("=== Seed Data Matched: 20 Unique Products Inserted! ===");
print("==========================================================");
