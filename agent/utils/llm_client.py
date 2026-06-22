"""LLM集成模块 - 用于自动代码生成和修复"""
import json
import os
from typing import Dict, List, Optional
from pathlib import Path

try:
    import requests
    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False


class LLMClient:
    """LLM客户端 - 支持OpenAI兼容API"""
    
    def __init__(self, api_url: str = None, api_key: str = None, model: str = None):
        self.api_url = api_url or os.getenv("LLM_API_URL", "https://api.openai.com/v1")
        self.api_key = api_key or os.getenv("LLM_API_KEY", "")
        self.model = model or os.getenv("LLM_MODEL", "gpt-4")
        self.timeout = 120
    
    def is_configured(self) -> bool:
        """检查LLM是否已配置"""
        return bool(self.api_key)
    
    def chat(self, messages: List[Dict], temperature: float = 0.7) -> Optional[str]:
        """发送聊天请求"""
        if not self.is_configured():
            return None
        
        if not HAS_REQUESTS:
            print("  [!] 需要安装requests库: pip install requests")
            return None
        
        try:
            headers = {
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json"
            }
            
            payload = {
                "model": self.model,
                "messages": messages,
                "temperature": temperature,
                "max_tokens": 4096
            }
            
            response = requests.post(
                f"{self.api_url}/chat/completions",
                headers=headers,
                json=payload,
                timeout=self.timeout
            )
            
            if response.status_code == 200:
                data = response.json()
                return data["choices"][0]["message"]["content"]
            else:
                print(f"  [!] LLM API错误: {response.status_code}")
                return None
                
        except Exception as e:
            print(f"  [!] LLM请求失败: {e}")
            return None
    
    def generate_code(self, prompt: str, context: str = "") -> Optional[str]:
        """生成代码"""
        messages = [
            {
                "role": "system",
                "content": """你是一个专业的Go/JavaScript/Vue开发者。
请根据要求生成高质量的代码。
只返回代码，不要解释。代码要完整可运行。"""
            }
        ]
        
        if context:
            messages.append({
                "role": "user",
                "content": f"项目上下文:\n{context}"
            })
        
        messages.append({
            "role": "user",
            "content": prompt
        })
        
        return self.chat(messages, temperature=0.3)
    
    def generate_test(self, source_code: str, language: str, module_name: str) -> Optional[str]:
        """为代码生成测试"""
        prompts = {
            "go": f"""为以下Go代码生成完整的单元测试。

要求:
1. 使用标准testing包
2. 测试所有导出函数
3. 包含正常和错误情况
4. 使用表驱动测试
5. 文件名: {module_name}_test.go

源代码:
{source_code}

只返回纯测试代码，不要包含任何Markdown标记或解释。""",
            
            "javascript": f"""为以下JavaScript代码生成完整的单元测试。

要求:
1. 使用Vitest框架
2. describe/it结构
3. 测试所有方法
4. Mock外部依赖
5. 文件名: {module_name}.test.js

源代码:
{source_code}

只返回纯测试代码，不要包含任何Markdown标记或解释。""",
            
            "vue": f"""为以下Vue组件/Store生成完整的单元测试。

要求:
1. 使用Vitest + @vue/test-utils
2. 测试所有props/events
3. 测试用户交互
4. Mock外部依赖
5. 文件名: {module_name}.test.js

源代码:
{source_code}

只返回纯测试代码，不要包含任何Markdown标记或解释。"""
        }
        
        prompt = prompts.get(language, prompts["javascript"])
        return self.generate_code(prompt)
    
    def fix_code(self, code: str, error: str, language: str) -> Optional[str]:
        """根据错误修复代码"""
        prompt = f"""修复以下代码中的错误。

错误信息:
{error}

代码:
```{language}
{code}
```

返回修复后的完整代码，不要解释。"""
        
        return self.generate_code(prompt)
    
    def analyze_and_improve(self, code: str, language: str) -> Optional[str]:
        """分析并改进代码"""
        prompt = f"""分析以下代码并提供改进建议，然后给出改进后的代码。

代码:
```{language}
{code}
```

要求:
1. 指出问题
2. 提供改进后的完整代码
3. 保持功能不变

返回格式:
## 问题
- 问题1
- 问题2

## 改进后的代码
(完整代码)"""
        
        return self.chat([{"role": "user", "content": prompt}])
    
    def generate_feature(self, task_description: str, requirements: List[str], 
                         language: str, context: str = "", file_type: str = "module") -> Optional[str]:
        """生成功能代码"""
        req_text = "\n".join(f"- {r}" for r in requirements)
        
        prompts = {
            "go": f"""根据以下需求生成Go语言功能代码。

任务描述: {task_description}

需求:
{req_text}

项目上下文:
{context}

要求:
1. 使用标准Go风格
2. 包含必要的imports
3. 添加适当的错误处理
4. 使用GORM进行数据库操作（如需要）
5. 遵循项目现有的代码风格
6. 代码完整可运行

只返回纯代码，不要包含Markdown标记或解释。""",
            
            "javascript": f"""根据以下需求生成功能代码。

任务描述: {task_description}

需求:
{req_text}

项目上下文:
{context}

要求:
1. 使用ES6+语法
2. 遵循项目现有的代码风格
3. 包含必要的导入
4. 添加适当的错误处理
5. 代码完整可运行

只返回纯代码，不要包含Markdown标记或解释。""",
            
            "vue": f"""根据以下需求生成Vue组件或Store代码。

任务描述: {task_description}

需求:
{req_text}

项目上下文:
{context}

要求:
1. 使用Composition API
2. 使用Pinia进行状态管理（如是Store）
3. 遵循项目现有的代码风格
4. 包含必要的导入
5. 代码完整可运行

只返回纯代码，不要包含Markdown标记或解释。"""
        }
        
        prompt = prompts.get(language, prompts["javascript"])
        return self.generate_code(prompt, context)
