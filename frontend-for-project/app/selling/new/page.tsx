"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Navbar } from "@/components/navbar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { fetchWithAuth } from "@/lib/api-client"
import { ChevronLeft } from "lucide-react"

export default function NewListingPage() {
  const router = useRouter()
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    categoryId: "",
    startingPrice: "",
    buyNowPrice: "",
    stepPrice: "",
    endAt: "",
    autoExtend: true,
  })

  useEffect(() => {
    async function loadCategories() {
      const res = await fetchWithAuth("/categories")
      const data = await res.json()
      setCategories(data.categories || [])
    }
    loadCategories()
  }, [])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const payload = {
        ...formData,
        categoryId: Number(formData.categoryId),
        startingPrice: Number(formData.startingPrice),
        buyNowPrice: formData.buyNowPrice ? Number(formData.buyNowPrice) : null,
        stepPrice: Number(formData.stepPrice),
        endAt: new Date(formData.endAt).toISOString().split(".")[0], // Match API format
      }

      const res = await fetchWithAuth("/products", {
        method: "POST",
        body: JSON.stringify(payload),
      })

      if (res.ok) {
        router.push("/selling")
      } else {
        const error = await res.json()
        alert(error.message || "Failed to create listing")
      }
    } catch (err) {
      alert("An error occurred")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-black pb-20">
      <Navbar />
      <main className="container mx-auto px-6 pt-32">
        <div className="max-w-3xl mx-auto flex flex-col gap-8">
          <Button
            variant="ghost"
            onClick={() => router.back()}
            className="w-fit text-zinc-400 hover:text-white -ml-4 rounded-full"
          >
            <ChevronLeft className="w-5 h-5 mr-1" />
            Back to Dashboard
          </Button>

          <div>
            <h1 className="text-4xl font-bold tracking-tighter mb-2">List New Product</h1>
            <p className="text-zinc-500">Enter the details for your premium auction item.</p>
          </div>

          <form onSubmit={handleSubmit} className="grid gap-12">
            <section className="flex flex-col gap-6">
              <h2 className="text-xl font-bold tracking-tight pb-2 border-b border-white/10">General Information</h2>
              <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Product Name</label>
                  <Input
                    required
                    className="h-12 bg-white/5 border-white/10 rounded-xl"
                    placeholder="e.g. 2024 Rolex Submariner"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  />
                </div>
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Category</label>
                  <Select onValueChange={(val) => setFormData({ ...formData, categoryId: val })}>
                    <SelectTrigger className="h-12 bg-white/5 border-white/10 rounded-xl">
                      <SelectValue placeholder="Select a category" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories.map((cat: any) => (
                        <SelectItem key={cat.id} value={cat.id.toString()}>
                          {cat.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Description</label>
                  <Textarea
                    required
                    className="min-h-[160px] bg-white/5 border-white/10 rounded-2xl p-4"
                    placeholder="Provide a detailed description..."
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  />
                </div>
              </div>
            </section>

            <section className="flex flex-col gap-6">
              <h2 className="text-xl font-bold tracking-tight pb-2 border-b border-white/10">Pricing & Timing</h2>
              <div className="grid md:grid-cols-2 gap-4">
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">
                    Starting Price ($)
                  </label>
                  <Input
                    type="number"
                    required
                    className="h-12 bg-white/5 border-white/10 rounded-xl"
                    value={formData.startingPrice}
                    onChange={(e) => setFormData({ ...formData, startingPrice: e.target.value })}
                  />
                </div>
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Buy Now Price ($)</label>
                  <Input
                    type="number"
                    className="h-12 bg-white/5 border-white/10 rounded-xl"
                    value={formData.buyNowPrice}
                    onChange={(e) => setFormData({ ...formData, buyNowPrice: e.target.value })}
                  />
                </div>
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">Step Price ($)</label>
                  <Input
                    type="number"
                    required
                    className="h-12 bg-white/5 border-white/10 rounded-xl"
                    value={formData.stepPrice}
                    onChange={(e) => setFormData({ ...formData, stepPrice: e.target.value })}
                  />
                </div>
                <div className="flex flex-col gap-2">
                  <label className="text-xs font-bold uppercase tracking-widest text-zinc-500">End Date & Time</label>
                  <Input
                    type="datetime-local"
                    required
                    className="h-12 bg-white/5 border-white/10 rounded-xl"
                    value={formData.endAt}
                    onChange={(e) => setFormData({ ...formData, endAt: e.target.value })}
                  />
                </div>
              </div>
              <div className="flex items-center justify-between p-6 rounded-2xl bg-white/5 border border-white/10">
                <div>
                  <p className="font-bold">Auto-Extend Auction</p>
                  <p className="text-sm text-zinc-500">Extend by 10m if a bid is placed in the last 5m.</p>
                </div>
                <Switch
                  checked={formData.autoExtend}
                  onCheckedChange={(checked) => setFormData({ ...formData, autoExtend: checked })}
                />
              </div>
            </section>

            <div className="pt-8 border-t border-white/10 flex justify-end gap-4">
              <Button
                type="button"
                variant="ghost"
                onClick={() => router.back()}
                className="h-12 px-8 rounded-full text-zinc-400 hover:text-white"
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={loading}
                className="h-12 px-12 bg-white text-black hover:bg-zinc-200 rounded-full font-bold"
              >
                {loading ? "Listing..." : "Publish Auction"}
              </Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  )
}
