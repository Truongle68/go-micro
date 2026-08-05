db = db.getSiblingDB("catalog-db");

db.createUser({
  user: "catalog_user",
  pwd: "catalog_password",
  roles: [
    {
      role: "readWrite",
      db: "catalog-db",
    },
  ],
});

db.createCollection("categories");
db.createCollection("products");

db.categories.createIndex({ slug: 1 }, { unique: true });

db.products.createIndex({ category_id: 1 });
db.products.createIndex({ "variants.sku": 1 }, { unique: true });
db.products.createIndex({ slug: 1 }, { unique: true });

print("=== Catalog Database initialized & indexes created ===");
