import type { FC } from 'react'
import type { ModelProvider } from '../declarations'
import { basePath } from '@/utils/var'
import { useLanguage } from '../hooks'
import { Openai } from '@/app/components/base/icons/src/vender/other'
import { AnthropicDark, AnthropicLight } from '@/app/components/base/icons/src/public/llm'
import { renderI18nObject } from '@/i18n'
import { Theme } from '@/types/app'
import cn from '@/utils/classnames'
import useTheme from '@/hooks/use-theme'

type ProviderIconProps = {
  provider: ModelProvider
  className?: string
}
const ProviderIcon: FC<ProviderIconProps> = ({
  provider,
  className,
}) => {
  const { theme } = useTheme()
  const language = useLanguage()

  if (provider.provider === 'khulnasoft/anthropic/anthropic') {
    return (
      <div className='mb-2 py-[7px]'>
        {theme === Theme.dark && <AnthropicLight className='h-2.5 w-[90px]' />}
        {theme === Theme.light && <AnthropicDark className='h-2.5 w-[90px]' />}
      </div>
    )
  }

  if (provider.provider === 'khulnasoft/openai/openai') {
    return (
      <div className='mb-2'>
        <Openai className='h-6 w-auto text-text-inverted-dimmed' />
      </div>
    )
  }

  return (
    <div className={cn('inline-flex items-center gap-2', className)}>
      <img
        alt='provider-icon'
        src={basePath + renderI18nObject(provider.icon_small, language)}
        className='h-6 w-6'
      />
      <div className='system-md-semibold text-text-primary'>
        {renderI18nObject(provider.label, language)}
      </div>
    </div>
  )
}

export default ProviderIcon
