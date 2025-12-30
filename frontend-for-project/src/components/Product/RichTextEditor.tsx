import { useState } from 'react';
import {
  Bold,
  Italic,
  List,
  ListOrdered,
  Link as LinkIcon,
} from 'lucide-react';

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  error?: string;
  disabled?: boolean;
}

export const RichTextEditor = ({
  value,
  onChange,
  placeholder = 'Nhập mô tả sản phẩm...',
  error,
  disabled = false,
}: RichTextEditorProps) => {
  const [showPreview, setShowPreview] = useState(false);

  const insertMarkdown = (before: string, after: string = '') => {
    const textarea = document.getElementById(
      'description-textarea'
    ) as HTMLTextAreaElement;
    if (!textarea) return;

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selectedText = value.substring(start, end);
    const newValue =
      value.substring(0, start) +
      before +
      selectedText +
      after +
      value.substring(end);

    onChange(newValue);

    // Restore cursor position
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(
        start + before.length,
        start + before.length + selectedText.length
      );
    }, 0);
  };

  const renderPreview = () => {
    // Simple markdown rendering (you can enhance this)
    let html = value;
    // Bold
    html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
    // Italic
    html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');
    // Headers
    html = html.replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-4 mb-2">$1</h3>');
    html = html.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-4 mb-2">$1</h2>');
    html = html.replace(/^# (.*$)/gim, '<h2 class="text-2xl font-bold mt-4 mb-2">$1</h2>');
    // Line breaks
    html = html.replace(/\n/g, '<br />');
    return { __html: html };
  };

  return (
    <div className="border border-gray-300 rounded-lg overflow-hidden">
      {/* Toolbar */}
      <div className="bg-gray-50 border-b border-gray-300 px-3 py-2 flex items-center gap-2 flex-wrap">
        <button
          type="button"
          onClick={() => insertMarkdown('**', '**')}
          disabled={disabled}
          className="p-2 hover:bg-gray-200 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Bold"
        >
          <Bold className="w-4 h-4" />
        </button>
        <button
          type="button"
          onClick={() => insertMarkdown('*', '*')}
          disabled={disabled}
          className="p-2 hover:bg-gray-200 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Italic"
        >
          <Italic className="w-4 h-4" />
        </button>
        <div className="w-px h-6 bg-gray-300 mx-1"></div>
        <button
          type="button"
          onClick={() => insertMarkdown('## ', '')}
          disabled={disabled}
          className="px-3 py-1 hover:bg-gray-200 rounded transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          title="Heading"
        >
          H
        </button>
        <button
          type="button"
          onClick={() => insertMarkdown('- ', '')}
          disabled={disabled}
          className="p-2 hover:bg-gray-200 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Bullet List"
        >
          <List className="w-4 h-4" />
        </button>
        <button
          type="button"
          onClick={() => insertMarkdown('1. ', '')}
          disabled={disabled}
          className="p-2 hover:bg-gray-200 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Numbered List"
        >
          <ListOrdered className="w-4 h-4" />
        </button>
        <button
          type="button"
          onClick={() => insertMarkdown('[', '](url)')}
          disabled={disabled}
          className="p-2 hover:bg-gray-200 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Link"
        >
          <LinkIcon className="w-4 h-4" />
        </button>
        <div className="w-px h-6 bg-gray-300 mx-1"></div>
        <button
          type="button"
          onClick={() => setShowPreview(!showPreview)}
          className="px-3 py-1 hover:bg-gray-200 rounded transition-colors text-sm font-medium"
        >
          {showPreview ? 'Chỉnh sửa' : 'Xem trước'}
        </button>
      </div>

      {/* Editor/Preview */}
      {showPreview ? (
        <div
          className="p-4 min-h-[250px] prose max-w-none"
          dangerouslySetInnerHTML={renderPreview()}
        />
      ) : (
        <textarea
          id="description-textarea"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          disabled={disabled}
          className={`w-full px-4 py-3 min-h-[250px] focus:outline-none resize-y ${
            disabled ? 'bg-gray-50 cursor-not-allowed' : ''
          }`}
        />
      )}

      {error && (
        <div className="px-4 py-2 bg-red-50 border-t border-red-200">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      {/* Markdown Guide */}
      <div className="px-4 py-2 bg-gray-50 border-t border-gray-200">
        <p className="text-xs text-gray-500">
          <strong>Hỗ trợ Markdown:</strong> **bold**, *italic*, ## heading, -
          list, [link](url)
        </p>
      </div>
    </div>
  );
};
