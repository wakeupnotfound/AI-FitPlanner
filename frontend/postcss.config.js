import autoprefixer from 'autoprefixer'
import pxtorem from 'postcss-pxtorem'

export default {
  plugins: [
    autoprefixer(),
    pxtorem({
      rootValue: 16, // 改为16，基于浏览器默认字体大小
      propList: ['*'],
      selectorBlackList: ['.norem'],
      exclude: /node_modules/i,
      minPixelValue: 2 // 小于2px的不转换
    })
  ]
}
