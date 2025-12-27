"use client"

import { useState, useEffect } from "react"
import { Navbar } from "@/components/navbar"
import { ProductCard } from "@/components/product-card"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Search } from "lucide-react"

export default function ProductsPage() {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchProducts() {
      try {
        const res = await fetch("http://localhost:8080/api/products/search?pageSize=12")
        const data = await res.json()
        if (data.success) {
          setProducts(data.data.content)
        }
      } catch (error) {
        console.error("Failed to fetch products:", error)
      } finally {
        setLoading(false)
      }
    }
    fetchProducts()
  }, [])

  return (
    <div className="min-h-screen pb-20">
      <Navbar />

      <main className="container mx-auto px-6 pt-32">
        <div className="flex flex-col gap-8">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div>
              <h1 className="text-4xl font-bold tracking-tighter mb-2">Live Auctions</h1>
              <p className="text-zinc-400">Discover and bid on the most exclusive products.</p>
            </div>

            <div className="flex items-center gap-3 w-full md:w-auto">
              <div className="relative flex-1 md:w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-500" />
                <Input
                  placeholder="Search products..."
                  className="pl-9 rounded-full bg-white/5 border-white/10 focus:border-white/20 h-10"
                />
              </div>
              <Select defaultValue="newest">
                <SelectTrigger className="w-[140px] rounded-full bg-white/5 border-white/10 h-10">
                  <SelectValue placeholder="Sort by" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="newest">Newest</SelectItem>
                  <SelectItem value="ending">Ending Soon</SelectItem>
                  <SelectItem value="price-asc">Price: Low to High</SelectItem>
                  <SelectItem value="price-desc">Price: High to Low</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {loading ? (
              Array(8)
                .fill(0)
                .map((_, i) => <div key={i} className="aspect-[3/4] rounded-xl bg-white/5 animate-pulse" />)
            ) : products.length > 0 ? (
              products.map((product: any) => <ProductCard key={product.id} product={product} />)
            ) : (
              <div className="col-span-full py-20 text-center">
                <p className="text-zinc-500">No products found. Please try a different search.</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
