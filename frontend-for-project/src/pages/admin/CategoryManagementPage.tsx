import { useState, useEffect, useCallback } from 'react';
import { categoryService } from '../../services/category.service';
import { Category } from '../../types';
import { useAuth } from '../../hooks/useAuth';
import { useUIStore } from '../../stores/ui.store';
import CategoryForm from '../../components/Admin/CategoryForm';
import CategoryList from '../../components/Admin/CategoryList';
import Loading from '../../components/Common/Loading';

export default function CategoryManagementPage() {
  const { user } = useAuth();
  const addToast = useUIStore((state) => state.addToast);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);

  const fetchCategories = useCallback(async () => {
    try {
      setLoading(true);
      const data = await categoryService.getAllCategories();
      setCategories(data);
    } catch (error) {
      addToast('error', 'Không thể tải danh sách danh mục');
      console.error(error);
    } finally {
      setLoading(false);
    }
  }, [addToast]);

  useEffect(() => {
    if (user?.userRole !== 'ROLE_ADMIN') {
      addToast('error', 'Bạn không có quyền truy cập trang này');
      window.location.href = '/';
      return;
    }
    fetchCategories();
  }, [user, addToast, fetchCategories]);

  const handleCreate = () => {
    setEditingCategory(null);
    setShowForm(true);
  };

  const handleEdit = (category: Category) => {
    setEditingCategory(category);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Bạn có chắc chắn muốn xóa danh mục này?')) {
      return;
    }

    try {
      await categoryService.deleteCategory(id);
      addToast('success', 'Xóa danh mục thành công');
      fetchCategories();
    } catch (error) {
      const errorMessage = error instanceof Error && 'response' in error 
        ? (error as unknown as { response?: { data?: { error?: string } } }).response?.data?.error 
        : 'Không thể xóa danh mục';
      addToast('error', errorMessage || 'Không thể xóa danh mục');
      console.error(error);
    }
  };

  const handleFormSubmit = async (data: Partial<Category>) => {
    try {
      if (editingCategory) {
        await categoryService.updateCategory(editingCategory.id, data);
        addToast('success', 'Cập nhật danh mục thành công');
      } else {
        await categoryService.createCategory(data as Category & { level: 1 | 2 });
        addToast('success', 'Tạo danh mục thành công');
      }
      setShowForm(false);
      setEditingCategory(null);
      fetchCategories();
    } catch (error) {
      const errorMessage = error instanceof Error && 'response' in error 
        ? (error as unknown as { response?: { data?: { error?: string } } }).response?.data?.error 
        : 'Có lỗi xảy ra';
      addToast('error', errorMessage || 'Có lỗi xảy ra');
      console.error(error);
    }
  };

  const handleFormCancel = () => {
    setShowForm(false);
    setEditingCategory(null);
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Quản Lý Danh Mục</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + Thêm Danh Mục
        </button>
      </div>

      {showForm && (
        <div className="mb-6">
          <CategoryForm
            category={editingCategory}
            categories={categories}
            onSubmit={handleFormSubmit}
            onCancel={handleFormCancel}
          />
        </div>
      )}

      <CategoryList
        categories={categories}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />
    </div>
  );
}
