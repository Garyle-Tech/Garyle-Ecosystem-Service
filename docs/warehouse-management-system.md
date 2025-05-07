# Warehouse Management System

WMS merupakan sistem perangkat lunak yang digunakan untuk mengelola operasi gudang, mulai dari penerimaan, penyimpanan, hingga pengiriman barang

Fitur - Fitur dari WMS yang akan di buat diantaranya :

- **Master Data**
| Entitas      | Atribut Utama                                                   | Endpoint RESTful                 | Method              | Keterangan                        |
| ------------ | --------------------------------------------------------------- | -------------------------------- | ------------------- | --------------------------------- |
| **Product**  | `id`, `sku`, `name`, `description`, `unit`, `weight`, `dimensi` | `/products`                      | GET                 | list semua produk                 |
|              |                                                                 | `/products/{id}`                 | GET                 | detail produk                     |
|              |                                                                 | `/products`                      | POST                | tambah produk baru                |
|              |                                                                 | `/products/{id}`                 | PUT                 | update data produk                |
|              |                                                                 | `/products/{id}`                 | DELETE              | hapus (atau soft-delete) produk   |
| **Location** | `id`, `code`, `zone`, `type` (rack/bin/area), `capacity`        | `/locations`                     | GET                 | list semua lokasi                 |
|              |                                                                 | `/locations/{id}`                | GET                 | detail lokasi                     |
|              |                                                                 | `/locations`                     | POST                | buat lokasi baru                  |
|              |                                                                 | `/locations/{id}`                | PATCH               | ubah kapabilitas atau tipe lokasi |
|              |                                                                 | `/locations/{id}`                | DELETE              | hapus lokasi                      |
| **Supplier** | `id`, `name`, `address`, `contact`                              | `/suppliers` `/suppliers/{id}`   | GET/POST/PUT/DELETE | CRUD supplier                     |
| **Customer** | `id`, `name`, `address`, `contact`                              | `/customers` `/customers/{id}`   | GET/POST/PUT/DELETE | CRUD pelanggan                    |
| **Category** | `id`, `name`, `parent_id`                                       | `/categories` `/categories/{id}` | GET/POST/PUT/DELETE | CRUD kategori                     |
| **Batch**    | `id`, `product_id`, `batch_no`, `expire_date`                   | `/batches` `/batches/{id}`       | GET/POST/PUT/DELETE | CRUD batch                        |

----
- **Inbound (Penerimaan Barang)**
| Entitas      | Atribut Utama                                                         | Endpoint                                   | Method              | Keterangan              |
| ------------ | --------------------------------------------------------------------- | ------------------------------------------ | ------------------- | ----------------------- |
| **PO**       | `id`, `supplier_id`, `order_date`, `lines: [{product_id, qty}]`       | `/purchase-orders` `/purchase-orders/{id}` | GET/POST/PUT/DELETE | CRUD PO                 |
| **ASN**      | `id`, `po_id`, `expected_date`, `lines`                               | `/asns` `/asns/{id}`                       | GET/POST/PUT/DELETE | CRUD ASN                |
| **Receipt**  | `id`, `asn_id`, `receive_date`, `lines: [{product_id, qty_received}]` | `/receipts` `/receipts/{id}`               | GET/POST            | buat & lihat penerimaan |
| **QCReport** | `id`, `receipt_id`, `status` (pass/reject), `notes`                   | `/qc-reports` `/qc-reports/{id}`           | GET/POST            | laporan QC penerimaan   |

---
- **Put Away**
| Entitas                | Atribut Utama                                                | Endpoint                               | Method       | Keterangan                          |
| ---------------------- | ------------------------------------------------------------ | -------------------------------------- | ------------ | ----------------------------------- |
| **PutawayTask**        | `id`, `receipt_id`, `from_location`, `to_location`, `status` | `/putaway-tasks` `/putaway-tasks/{id}` | GET/POST/PUT | CRUD & update status tugas put-away |
| **SlotRecommendation** | (di-generate, read-only)                                     | `/slot-recommendations`                | GET          | rekomendasi lokasi penempatan       |

---
- Inventory
| Entitas         | Atribut Utama                                                 | Endpoint                             | Method              | Keterangan             |
| --------------- | ------------------------------------------------------------- | ------------------------------------ | ------------------- | ---------------------- |
| **Inventory**   | `product_id`, `location_id`, `qty_on_hand`, `qty_available`   | `/inventory`                         | GET                 | lihat semua stok       |
|                 |                                                               | `/inventory/{product_id}`            | GET                 | stok per produk        |
| **CycleCount**  | `id`, `inventory_id`, `count_date`, `qty_counted`, `variance` | `/cycle-counts` `/cycle-counts/{id}` | GET/POST/PUT/DELETE | manajemen stock opname |
| **SerialTrack** | `serial_no`, `product_id`, `location_id`, `status`            | `/serials` `/serials/{serial_no}`    | GET/POST/PUT/DELETE | pelacakan serial       |


--- 
- **Order Fulfillment (Picking & Packing)**
| Entitas         | Atribut Utama                                                           | Endpoint                         | Method       | Keterangan                   |
| --------------- | ----------------------------------------------------------------------- | -------------------------------- | ------------ | ---------------------------- |
| **PickTask**    | `id`, `order_id`, `lines: [{product_id, qty}]`, `assigned_to`, `status` | `/pick-tasks` `/pick-tasks/{id}` | GET/POST/PUT | CRUD & update status picking |
| **PackTask**    | `id`, `pick_task_id`, `package_spec`, `status`                          | `/pack-tasks` `/pack-tasks/{id}` | GET/POST/PUT | CRUD & update status packing |
| **PickList**    | (read-only) gabungan pick-tasks                                         | `/orders/{id}/pick-list`         | GET          | generate daftar picking      |
| **PackingSlip** | (read-only PDF)                                                         | `/orders/{id}/packing-slip`      | GET          | cetak packing slip           |


--- 
- **Outbound**
| Entitas        | Atribut Utama                                                 | Endpoint                             | Method              | Keterangan                   |
| -------------- | ------------------------------------------------------------- | ------------------------------------ | ------------------- | ---------------------------- |
| **SalesOrder** | `id`, `customer_id`, `order_date`, `lines`, `status`          | `/sales-orders` `/sales-orders/{id}` | GET/POST/PUT/DELETE | CRUD order pelanggan         |
| **Shipment**   | `id`, `sales_order_id`, `ship_date`, `carrier`, `tracking_no` | `/shipments` `/shipments/{id}`       | GET/POST/PUT        | CRUD & konfirmasi kirim      |
| **Manifest**   | (read-only PDF)                                               | `/shipments/{id}/manifest`           | GET                 | generate manifest pengiriman |


---
- **Transfer Antar Lokasi**
| Entitas             | Atribut Utama                                       | Endpoint                       | Method       | Keterangan                      |
| ------------------- | --------------------------------------------------- | ------------------------------ | ------------ | ------------------------------- |
| **TransferRequest** | `id`, `from_loc_id`, `to_loc_id`, `lines`, `status` | `/transfers` `/transfers/{id}` | GET/POST/PUT | manage permintaan transfer stok |
| **TransferConfirm** | `id`, `transfer_id`, `confirmed_by`, `confirm_date` | `/transfers/{id}/confirm`      | POST         | konfirmasi pemindahan           |


---
- **Returns Management (RMA)**
| Entitas        | Atribut Utama                               | Endpoint              | Method              | Keterangan            |
| -------------- | ------------------------------------------- | --------------------- | ------------------- | --------------------- |
| **RMARequest** | `id`, `sales_order_id`, `reason`, `status`  | `/rmas` `/rmas/{id}`  | GET/POST/PUT/DELETE | manage RMA            |
| **RMAReceipt** | `id`, `rma_id`, `receive_date`, `qc_status` | `/rmas/{id}/receipts` | POST                | terima retur & QC     |
| **RMAPutaway** | `id`, `receipt_id`, `to_location`, `status` | `/rmas/{id}/putaway`  | POST                | restock setelah retur |


--- 
- **Task & Workforce**
| Entitas            | Atribut Utama                                            | Endpoint                           | Method              | Keterangan            |
| ------------------ | -------------------------------------------------------- | ---------------------------------- | ------------------- | --------------------- |
| **Worker**         | `id`, `name`, `role`, `shift`                            | `/workers` `/workers/{id}`         | GET/POST/PUT/DELETE | CRUD data operator    |
| **TaskAssignment** | `id`, `task_type`, `task_id`, `worker_id`, `assigned_at` | `/assignments` `/assignments/{id}` | GET/POST/PUT        | assign & update tugas |
| **Performance**    | (read-only, statistik)                                   | `/workers/{id}/performance`        | GET                 | metrik produktivitas  |


----
-  **Reporting & Analytics**
| Entitas / Endpoint         | Output / Keterangan                 | Method |
| -------------------------- | ----------------------------------- | ------ |
| `/reports/inventory-aging` | stock aging report                  | GET    |
| `/reports/throughput`      | throughput per periode              | GET    |
| `/reports/stock-turnover`  | inventory turnover rate             | GET    |
| `/reports/audit-trail`     | histori semua perubahan stok & data | GET    |

