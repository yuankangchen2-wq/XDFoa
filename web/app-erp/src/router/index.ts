import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/stocks' },
  { path: '/stocks', component: () => import('../views/StockList.vue') },
  { path: '/warehouses', component: () => import('../views/WarehouseList.vue') },
  { path: '/journals', component: () => import('../views/JournalList.vue') },
  { path: '/transfers', component: () => import('../views/TransferList.vue') },
  { path: '/stocktakes', component: () => import('../views/StocktakeList.vue') },
  { path: '/purchase/orders', component: () => import('../views/PurchaseOrderList.vue') },
  { path: '/purchase/receipts', component: () => import('../views/PurchaseReceiptList.vue') },
  { path: '/production/boms', component: () => import('../views/BomList.vue') },
  { path: '/production/work-orders', component: () => import('../views/WorkOrderList.vue') },
  { path: '/finance/receivables', component: () => import('../views/ReceivableList.vue') },
  { path: '/finance/payables', component: () => import('../views/PayableList.vue') },
  { path: '/finance/payments', component: () => import('../views/PaymentList.vue') },
  { path: '/cost/centers', component: () => import('../views/CostCenterList.vue') },
  { path: '/cost/products', component: () => import('../views/ProductCostList.vue') },
  { path: '/cost/records', component: () => import('../views/CostRecordList.vue') }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
