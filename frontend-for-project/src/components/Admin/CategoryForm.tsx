import { useState, useEffect } from 'react';
import { Category } from '../../types';

interface CategoryFormProps {
  category: Category | null;
  categories: Category[];
  onSubmit: (data: Partial<Category>) => void;
  onCancel: () => void;
}

export default function CategoryForm({ category, categories, onSubmit, onCancel }: CategoryFormProps) {
  const [formData, setFormData] = useState({
    name: '',
    slug: '',
    description: '',
    parent_id: undefined as number | undefined,
    level: 1 as 1 | 2,
    display_order: 0,
    is_active: true,
  });

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name,
        slug: category.slug,
        description: category.description || '',
        parent_id: category.parent_id,
        level: category.level,
        display_order: category.display_order,
        is_active: category.is_active,
      });
    }
  }, [category]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const checked = (e.target as HTMLInputElement).checked;
    
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : 
               name === 'parent_id' && value === '' ? undefined :
               name === 'parent_id' || name === 'display_order' ? Number(value) :
               value
    }));
  };

  const handleLevelChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const level = Number(e.target.value) as 1 | 2;
    setFormData(prev => ({
      ...prev,
      level,
      parent_id: level === 1 ? undefined : prev.parent_id
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Generate slug from name if not provided
    const slug = formData.slug || formData.name.toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/đ/g, 'd')
      .replace(/Đ/g, 'D')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');

    onSubmit({
      ...formData,
      slug,
      parent_id: formData.level === 1 ? undefined : formData.parent_id
    });
  };

  const parentCategories = categories.filter(c => c.level === 1);

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">
        {category ? 'Cập Nhật Danh Mục' : 'Tạo Danh Mục Mới'}
      </h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Name */}
          <div>
            <label className="block text-sm font-medium mb-1">
              Tên Danh Mục <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              required
              className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Điện tử, Thời trang..."
            />
          </div>

          {/* Slug */}
          <div>
            <label className="block text-sm font-medium mb-1">
              Slug (URL-friendly)
            </label>
            <input
              type="text"
              name="slug"
              value={formData.slug}
              onChange={handleChange}
              className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="dien-tu, thoi-trang..."
            />
            <p className="text-xs text-gray-500 mt-1">
              Để trống để tự động tạo từ tên danh mục
            </p>
          </div>

          {/* Level */}
          <div>
            <label className="block text-sm font-medium mb-1">
              Cấp <span className="text-red-500">*</span>
            </label>
            <select
              name="level"
              value={formData.level}
              onChange={handleLevelChange}
              className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              disabled={!!category} // Cannot change level when editing
            >
              <option value={1}>Cấp 1 (Danh mục chính)</option>
              <option value={2}>Cấp 2 (Danh mục con)</option>
            </select>
          </div>

          {/* Parent Category - Only for Level 2 */}
          {formData.level === 2 && (
            <div>
              <label className="block text-sm font-medium mb-1">
                Danh Mục Cha <span className="text-red-500">*</span>
              </label>
              <select
                name="parent_id"
                value={formData.parent_id || ''}
                onChange={handleChange}
                required={formData.level === 2}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Chọn danh mục cha</option>
                {parentCategories.map(cat => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          {/* Display Order */}
          <div>
            <label className="block text-sm font-medium mb-1">
              Thứ Tự Hiển Thị
            </label>
            <input
              type="number"
              name="display_order"
              value={formData.display_order}
              onChange={handleChange}
              min={0}
              className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Active Status */}
          <div className="flex items-center">
            <input
              type="checkbox"
              name="is_active"
              checked={formData.is_active}
              onChange={handleChange}
              id="is_active"
              className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label htmlFor="is_active" className="ml-2 text-sm font-medium">
              Kích hoạt
            </label>
          </div>
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium mb-1">
            Mô Tả
          </label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows={3}
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Mô tả ngắn về danh mục này..."
          />
        </div>

        {/* Action Buttons */}
        <div className="flex gap-3 pt-4">
          <button
            type="submit"
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            {category ? 'Cập Nhật' : 'Tạo Mới'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
          >
            Hủy
          </button>
        </div>
      </form>
    </div>
  );
}
