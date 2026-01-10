import { Check } from 'lucide-react';
import { OrderDetail } from '../../types';

interface OrderWizardProps {
  order: OrderDetail;
  isBuyer: boolean;
}

interface Step {
  key: string;
  label: string;
  description: string;
  statuses: OrderDetail['status'][];
}

const OrderWizard = ({ order, isBuyer }: OrderWizardProps) => {
  // Define workflow steps for buyer and seller
  const buyerSteps: Step[] = [
    {
      key: 'payment',
      label: 'Thanh toán',
      description: 'Hoàn tất thanh toán đơn hàng',
      statuses: ['PENDING_PAYMENT']
    },
    {
      key: 'address',
      label: 'Địa chỉ nhận hàng',
      description: 'Cung cấp địa chỉ giao hàng',
      statuses: ['PAID']
    },
    {
      key: 'shipping',
      label: 'Vận chuyển',
      description: 'Đợi người bán gửi hàng',
      statuses: ['ADDRESS_PROVIDED', 'SHIPPING']
    },
    {
      key: 'delivery',
      label: 'Nhận hàng',
      description: 'Xác nhận đã nhận hàng',
      statuses: ['SHIPPING']
    },
    {
      key: 'complete',
      label: 'Hoàn thành',
      description: 'Giao dịch thành công',
      statuses: ['DELIVERED', 'COMPLETED']
    }
  ];

  const sellerSteps: Step[] = [
    {
      key: 'payment',
      label: 'Chờ thanh toán',
      description: 'Đợi người mua thanh toán',
      statuses: ['PENDING_PAYMENT']
    },
    {
      key: 'address',
      label: 'Chờ địa chỉ',
      description: 'Đợi người mua cung cấp địa chỉ',
      statuses: ['PAID']
    },
    {
      key: 'shipping',
      label: 'Chuẩn bị hàng',
      description: 'Đóng gói và gửi hàng',
      statuses: ['ADDRESS_PROVIDED']
    },
    {
      key: 'delivery',
      label: 'Đang giao',
      description: 'Đơn hàng đang được vận chuyển',
      statuses: ['SHIPPING']
    },
    {
      key: 'complete',
      label: 'Hoàn thành',
      description: 'Giao dịch thành công',
      statuses: ['DELIVERED', 'COMPLETED']
    }
  ];

  const steps = isBuyer ? buyerSteps : sellerSteps;

  // Determine current step index based on order status
  const getCurrentStepIndex = () => {
    if (order.status === 'CANCELLED') return -1;
    
    for (let i = steps.length - 1; i >= 0; i--) {
      if (steps[i].statuses.includes(order.status)) {
        return i;
      }
    }
    return 0;
  };

  const currentStepIndex = getCurrentStepIndex();

  // Check if step is completed
  const isStepCompleted = (index: number) => {
    if (order.status === 'CANCELLED') return false;
    return index < currentStepIndex || order.status === 'COMPLETED';
  };

  // Check if step is current
  const isStepCurrent = (index: number) => {
    if (order.status === 'CANCELLED') return false;
    return index === currentStepIndex && order.status !== 'COMPLETED';
  };

  // Handle cancelled status
  if (order.status === 'CANCELLED') {
    return (
      <div className="bg-red-50 border-2 border-red-200 rounded-lg p-6">
        <div className="flex items-center justify-center">
          <div className="text-center">
            <div className="w-16 h-16 bg-red-500 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-white text-3xl">✕</span>
            </div>
            <h3 className="text-xl font-bold text-red-800 mb-2">Đơn hàng đã bị hủy</h3>
            {order.cancel_reason && (
              <p className="text-red-600">Lý do: {order.cancel_reason}</p>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">
          Quy trình đơn hàng
        </h2>
        <div className="text-sm text-gray-600">
          {order.status === 'COMPLETED' ? (
            <span className="text-green-600 font-medium">✓ Đã hoàn thành</span>
          ) : (
            <span>
              Bước {currentStepIndex + 1}/{steps.length}
            </span>
          )}
        </div>
      </div>
      
      {/* Desktop View - Horizontal */}
      <div className="hidden md:block">
        <div className="relative">
          {/* Progress Line */}
          <div className="absolute top-8 left-0 right-0 h-1 bg-gray-200">
            <div 
              className="h-full bg-blue-600 transition-all duration-500"
              style={{ 
                width: order.status === 'COMPLETED' 
                  ? '100%' 
                  : `${(currentStepIndex / (steps.length - 1)) * 100}%` 
              }}
            />
          </div>

          {/* Steps */}
          <div className="relative flex justify-between">
            {steps.map((step, index) => {
              const completed = isStepCompleted(index);
              const current = isStepCurrent(index);
              
              return (
                <div key={step.key} className="flex flex-col items-center" style={{ width: `${100 / steps.length}%` }}>
                  {/* Circle */}
                  <div
                    className={`w-16 h-16 rounded-full flex items-center justify-center border-4 transition-all duration-300 ${
                      completed
                        ? 'bg-blue-600 border-blue-600 text-white'
                        : current
                        ? 'bg-white border-blue-600 text-blue-600 ring-4 ring-blue-100'
                        : 'bg-white border-gray-300 text-gray-400'
                    }`}
                  >
                    {completed ? (
                      <Check className="w-8 h-8" />
                    ) : (
                      <span className="text-xl font-bold">{index + 1}</span>
                    )}
                  </div>
                  
                  {/* Label */}
                  <div className="mt-4 text-center max-w-[140px]">
                    <p className={`font-semibold mb-1 ${
                      completed || current ? 'text-gray-900' : 'text-gray-500'
                    }`}>
                      {step.label}
                    </p>
                    <p className={`text-xs ${
                      completed || current ? 'text-gray-600' : 'text-gray-400'
                    }`}>
                      {step.description}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Mobile View - Vertical */}
      <div className="md:hidden space-y-4">
        {steps.map((step, index) => {
          const completed = isStepCompleted(index);
          const current = isStepCurrent(index);
          
          return (
            <div key={step.key} className="flex gap-4">
              {/* Left side - Circle and Line */}
              <div className="flex flex-col items-center">
                <div
                  className={`w-12 h-12 rounded-full flex items-center justify-center border-4 transition-all duration-300 flex-shrink-0 ${
                    completed
                      ? 'bg-blue-600 border-blue-600 text-white'
                      : current
                      ? 'bg-white border-blue-600 text-blue-600 ring-4 ring-blue-100'
                      : 'bg-white border-gray-300 text-gray-400'
                  }`}
                >
                  {completed ? (
                    <Check className="w-6 h-6" />
                  ) : (
                    <span className="text-lg font-bold">{index + 1}</span>
                  )}
                </div>
                {index < steps.length - 1 && (
                  <div className={`w-1 flex-1 min-h-[40px] ${
                    completed ? 'bg-blue-600' : 'bg-gray-200'
                  }`} />
                )}
              </div>
              
              {/* Right side - Content */}
              <div className="flex-1 pb-6">
                <p className={`font-semibold mb-1 ${
                  completed || current ? 'text-gray-900' : 'text-gray-500'
                }`}>
                  {step.label}
                </p>
                <p className={`text-sm ${
                  completed || current ? 'text-gray-600' : 'text-gray-400'
                }`}>
                  {step.description}
                </p>
                {current && (
                  <div className="mt-2">
                    <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      <span className="w-2 h-2 bg-blue-600 rounded-full mr-2 animate-pulse"></span>
                      Đang thực hiện
                    </span>
                  </div>
                )}
                {completed && index === currentStepIndex - 1 && (
                  <div className="mt-2">
                    <span className="inline-flex items-center text-xs text-green-600">
                      ✓ Hoàn thành
                    </span>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default OrderWizard;
