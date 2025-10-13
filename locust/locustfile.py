# from locust import FastHttpUser, task, between
# import random

# class ProductUser(FastHttpUser):
#     # wait_time = between(1, 3)  # Simulate user think time (1–3s between requests)

#     @task(3)
#     def get_product(self):
#         # Pick a valid preloaded product ID (1–5)
#         product_id = random.choice([1, 2, 3, 4, 5])
#         self.client.get(f"/products/{product_id}", name="/products/:id")

#     @task(1)
#     def post_product_details(self):
#         # Randomly pick a product ID and modify it slightly
#         product_id = random.choice([1, 2, 3, 4, 5])
#         payload = {
#             "product_id": product_id,
#             "sku": f"SKU-{product_id}",
#             "manufacturer": "AcmeTest",
#             "category_id": 99,
#             "weight": 500 + product_id,
#             "some_other_id": 1000 + product_id
#         }
#         self.client.post(f"/products/{product_id}/details", json=payload, name="/products/:id/details")
# from locust import HttpUser, task, between
# import random

# class ProductUser(HttpUser):
#     # Simulate user "think time" between requests (1–3 seconds)
#     # wait_time = between(1, 3)

#     @task(3)
#     def get_product(self):
#         # Simulate browsing existing products (most common operation)
#         product_id = random.choice([1, 2, 3, 4, 5])
#         self.client.get(f"/products/{product_id}", name="/products/:id")

#     @task(1)
#     def post_product_details(self):
#         # Simulate admin updating product info occasionally
#         product_id = random.choice([1, 2, 3, 4, 5])
#         payload = {
#             "product_id": product_id,
#             "sku": f"SKU-{product_id}",
#             "manufacturer": "AcmeCorp",
#             "category_id": 99,
#             "weight": 500 + product_id,
#             "some_other_id": 2000 + product_id
#         }
#         self.client.post(f"/products/{product_id}/details", json=payload, name="/products/:id/details")


from locust import task, between
from locust.contrib.fasthttp import FastHttpUser
import random

# Common query terms that will always yield matches
QUERIES = [
    "electronics", "books", "home", "toys", "sports",
    "clothing", "beauty", "garden", "office", "grocery",
    "alpha", "beta", "gamma", "delta", "omega", ""
]

class SearchUser(FastHttpUser):
    # Minimal think time between requests
    network_timeout = 10.0
    connection_timeout = 10.0

    @task(10)
    def search_products(self):
        q = random.choice(QUERIES)
        with self.client.get(f"/products/search?q={q}", name="/products/search", catch_response=True) as resp:
            if resp.status_code != 200:
                resp.failure(f"HTTP {resp.status_code}")
            else:
                # Validate schema and constant-time behavior
                try:
                    data = resp.json()
                    if "total_found" not in data or "products" not in data:
                        resp.failure("Missing expected fields")
                    elif len(data["products"]) > 20:
                        resp.failure("More than 20 products returned")
                    else:
                        resp.success()
                except Exception:
                    resp.failure("Invalid JSON")
