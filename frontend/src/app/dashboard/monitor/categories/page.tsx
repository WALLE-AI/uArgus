import { FolderTree } from "lucide-react";

const CATEGORIES = [
  { name: "科技", count: 842, color: "bg-accent-blue" },
  { name: "金融", count: 634, color: "bg-accent-green" },
  { name: "地缘政治", count: 412, color: "bg-accent-purple" },
  { name: "能源", count: 287, color: "bg-accent-orange" },
  { name: "医疗健康", count: 195, color: "bg-red-400" },
  { name: "气候环境", count: 156, color: "bg-emerald-400" },
  { name: "学术研究", count: 321, color: "bg-indigo-400" },
  { name: "加密货币", count: 178, color: "bg-yellow-500" },
];

export default function CategoriesPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">领域分类</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {CATEGORIES.map((cat) => (
          <div key={cat.name} className="rounded-xl border border-border bg-white p-5 hover:shadow-md transition-shadow cursor-pointer">
            <div className="flex items-center gap-3">
              <div className={`rounded-lg ${cat.color}/10 p-2`}>
                <FolderTree className={`h-4 w-4 ${cat.color.replace("bg-", "text-")}`} />
              </div>
              <div>
                <div className="text-sm font-semibold">{cat.name}</div>
                <div className="text-xs text-muted-foreground">{cat.count} 数据源</div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
