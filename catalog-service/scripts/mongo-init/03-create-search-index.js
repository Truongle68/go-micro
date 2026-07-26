db = db.getSiblingDB("catalog-db");

const existingIndexes = db.products.getSearchIndexes({
  name: "products_fts_idx",
});

if (existingIndexes.length === 0) {
  db.products.createSearchIndex("products_fts_idx", {
    mappings: {
      dynamic: false,
      fields: {
        name_en: [
          { type: "string", analyzer: "lucene.standard" },
          { type: "autocomplete", analyzer: "lucene.standard" },
        ],
        name_vi: [
          { type: "string", analyzer: "lucene.standard" },
          { type: "autocomplete", analyzer: "lucene.standard" },
        ],
        description_en: { type: "string", analyzer: "lucene.standard" },
        description_vi: { type: "string", analyzer: "lucene.standard" },
        sku: { type: "string", analyzer: "lucene.keyword" },
        category_id: { type: "token" },
        is_active: { type: "boolean" },
        base_price: { type: "number" },
      },
    },
  });
  print("=== Search Index products_fts_idx created ===");
} else {
  print("=== Search Index products_fts_idx already exists ===");
}
