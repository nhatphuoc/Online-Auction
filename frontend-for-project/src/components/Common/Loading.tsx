export const Skeleton = ({ className = '' }: { className?: string }) => (
  <div className={`bg-gray-200 rounded animate-pulse ${className}`} />
);

export const ProductSkeleton = () => (
  <div className="bg-white rounded-lg shadow-sm overflow-hidden">
    <Skeleton className="h-48 w-full" />
    <div className="p-4 space-y-3">
      <Skeleton className="h-4 w-3/4" />
      <Skeleton className="h-4 w-1/2" />
      <Skeleton className="h-6 w-1/3" />
    </div>
  </div>
);

export const LoadingSpinner = () => (
  <div className="flex items-center justify-center py-12">
    <div className="w-8 h-8 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin" />
  </div>
);

export default LoadingSpinner;
