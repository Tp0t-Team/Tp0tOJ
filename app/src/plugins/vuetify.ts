import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import { VuetifyPreset } from 'vuetify/types/presets'
import colors from 'vuetify/es5/util/colors'

Vue.use(Vuetify)

const opt: VuetifyPreset = {
    icons: {
        iconfont: 'md'
    },
    theme: {
        dark: true,
        themes: {
            dark: {
                primary: colors.orange.darken2
            },
            light: {
                primary: colors.orange.darken2
            }
        }
    }
}

export default new Vuetify(opt)
