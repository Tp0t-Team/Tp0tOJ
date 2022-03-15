import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'
import { VuetifyPreset } from 'vuetify/types/services/presets'
import colors from 'vuetify/lib/util/colors'

Vue.use(Vuetify)

const opt = {
    icons: {
        iconfont: 'md',
        values: {}
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

export default new Vuetify(opt as any)
