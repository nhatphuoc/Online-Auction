import Link from "next/link"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Clock, Gavel } from "lucide-react"

interface ProductCardProps {
  product: {
    id: number
    name: string
    thumbnailUrl: string
    currentPrice: number
    buyNowPrice?: number
    bidCount: number
    endAt: string
    categoryName: string
    isNew?: boolean
  }
}

export function ProductCard({ product }: ProductCardProps) {
  const timeLeft = new Date(product.endAt).getTime() - Date.now()
  const days = Math.floor(timeLeft / (1000 * 60 * 60 * 24))
  const hours = Math.floor((timeLeft / (1000 * 60 * 60)) % 24)

  return (
    <Card className="overflow-hidden bg-white/5 border-white/10 hover:border-white/20 transition-all group">
      <div className="aspect-[4/3] relative overflow-hidden">
        <img
          src={product.thumbnailUrl || "/placeholder.svg?height=300&width=400&query=product"}
          alt={product.name}
          className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-500"
        />
        {product.isNew && (
          <Badge className="absolute top-3 left-3 bg-white text-black font-bold rounded-full">NEW</Badge>
        )}
      </div>
      <CardContent className="p-4 flex flex-col gap-2">
        <div className="flex justify-between items-start">
          <p className="text-zinc-500 text-xs font-medium uppercase tracking-wider">{product.categoryName}</p>
          <div className="flex items-center gap-1 text-zinc-400 text-xs">
            <Clock className="w-3 h-3" />
            <span>{days > 0 ? `${days}d ${hours}h` : `${hours}h left`}</span>
          </div>
        </div>
        <h3 className="font-bold text-lg leading-tight line-clamp-2 min-h-[2.5rem]">{product.name}</h3>
        <div className="flex flex-col gap-1 mt-2">
          <div className="flex justify-between items-baseline">
            <span className="text-zinc-400 text-sm">Current Bid</span>
            <span className="text-white font-bold text-xl">${product.currentPrice.toLocaleString()}</span>
          </div>
          {product.buyNowPrice && (
            <div className="flex justify-between items-baseline opacity-60">
              <span className="text-zinc-400 text-xs">Buy Now</span>
              <span className="text-white font-medium text-sm">${product.buyNowPrice.toLocaleString()}</span>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="p-4 pt-0 flex gap-2">
        <Link href={`/products/${product.id}`} className="flex-1">
          <Button className="w-full bg-white text-black hover:bg-zinc-200 rounded-full font-bold">
            <Gavel className="w-4 h-4 mr-2" />
            Place Bid
          </Button>
        </Link>
        <Badge variant="outline" className="h-9 px-3 rounded-full border-white/10 text-zinc-400 font-medium">
          {product.bidCount} bids
        </Badge>
      </CardFooter>
    </Card>
  )
}
