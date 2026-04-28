export default function SettingsPage() {
  return (
    <div className="space-y-6 max-w-2xl">
      <h1 className="text-2xl font-bold">用户设置</h1>

      <div className="rounded-xl border border-border bg-white p-6 space-y-6">
        <div>
          <h3 className="text-sm font-semibold mb-3">个人信息</h3>
          <div className="space-y-3">
            <div>
              <label className="text-xs text-muted-foreground">用户名</label>
              <input defaultValue="User" className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none focus:border-accent-blue" />
            </div>
            <div>
              <label className="text-xs text-muted-foreground">邮箱</label>
              <input defaultValue="user@example.com" className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none focus:border-accent-blue" />
            </div>
          </div>
        </div>

        <div className="border-t border-border pt-6">
          <h3 className="text-sm font-semibold mb-3">偏好设置</h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-sm">语言</div>
                <div className="text-xs text-muted-foreground">选择界面显示语言</div>
              </div>
              <select className="rounded-lg border border-border px-3 py-1.5 text-sm outline-none">
                <option>中文</option>
                <option>English</option>
              </select>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <div className="text-sm">主题</div>
                <div className="text-xs text-muted-foreground">选择深色或浅色模式</div>
              </div>
              <select className="rounded-lg border border-border px-3 py-1.5 text-sm outline-none">
                <option>浅色</option>
                <option>深色</option>
                <option>跟随系统</option>
              </select>
            </div>
          </div>
        </div>

        <div className="border-t border-border pt-6">
          <h3 className="text-sm font-semibold mb-3">API Key 管理</h3>
          <div className="rounded-lg border border-border p-3 flex items-center justify-between">
            <div>
              <div className="text-sm font-mono">sk-****...****3f2a</div>
              <div className="text-xs text-muted-foreground">创建于 2024-01-15</div>
            </div>
            <button className="text-xs text-red-500 hover:text-red-600">撤销</button>
          </div>
          <button className="mt-3 rounded-lg border border-border px-3 py-1.5 text-xs text-muted-foreground hover:bg-muted transition-colors">
            + 生成新 Key
          </button>
        </div>
      </div>
    </div>
  );
}
