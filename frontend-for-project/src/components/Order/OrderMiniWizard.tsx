import { Check } from 'lucide-react';
import { Order } from '../../types';

interface OrderMiniWizardProps {
  status: Order['status'];
  isBuyer: boolean;
}

const OrderMiniWizard = ({ status, isBuyer }: OrderMiniWizardProps) => {
  const buyerSteps = [
    { key: 'PENDING_PAYMENT', label: 'Thanh toán' },
    { key: 'PAID', label: 'Địa chỉ' },
    { key: 'ADDRESS_PROVIDED', label: 'Chờ gửi' },
    { key: 'SHIPPING', label: 'Vận chuyển' },
    { key: 'DELIVERED', label: 'Hoàn thành' },
  ];

  const sellerSteps = [
    { key: 'PENDING_PAYMENT', label: 'Chờ TT' },
    { key: 'PAID', label: 'Chờ địa chỉ' },
    { key: 'ADDRESS_PROVIDED', label: 'Gửi hàng' },
    { key: 'SHIPPING', label: 'Đang giao' },
    { key: 'DELIVERED', label: 'Hoàn thành' },
  ];

  const steps = isBuyer ? buyerSteps : sellerSteps;

  // Map status to step index
  const statusToIndex: Record<Order['status'], number> = {
    PENDING_PAYMENT: 0,
    PAID: 1,
    ADDRESS_PROVIDED: 2,
    SHIPPING: 3,
    DELIVERED: 4,
    COMPLETED: 4,
    CANCELLED: -1,
  };

  const currentIndex = statusToIndex[status];

  if (status === 'CANCELLED') {
    return (
      <div className="flex items-center gap-2 text-red-600">
        <span className="text-xs font-medium">✕ Đã hủy</span>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-1">
      {steps.map((step, index) => {
        const isCompleted = index < currentIndex || status === 'COMPLETED';
        const isCurrent = index === currentIndex && status !== 'COMPLETED';
        
        return (
          <div key={step.key} className="flex items-center">
            <div className="flex flex-col items-center">
              <div
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium transition-colors ${
                  isCompleted
                    ? 'bg-blue-600 text-white'
                    : isCurrent
                    ? 'bg-blue-100 text-blue-600 border-2 border-blue-600'
                    : 'bg-gray-200 text-gray-500'
                }`}
              >
                {isCompleted ? <Check className="w-3 h-3" /> : index + 1}
              </div>
              <span className={`text-[10px] mt-1 whitespace-nowrap ${
                isCompleted || isCurrent ? 'text-gray-700 font-medium' : 'text-gray-400'
              }`}>
                {step.label}
              </span>
            </div>
            {index < steps.length - 1 && (
              <div className={`w-6 h-0.5 mx-0.5 mb-4 ${
                isCompleted ? 'bg-blue-600' : 'bg-gray-200'
              }`} />
            )}
          </div>
        );
      })}
    </div>
  );
};

export default OrderMiniWizard;
