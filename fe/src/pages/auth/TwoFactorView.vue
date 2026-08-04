<script setup lang="ts">
import { FormInput, UiButton, UiCard } from '@/components/utils'
import { ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { required, numeric, minLength, maxLength } from '@vuelidate/validators'
import swal from '@/plugins/swal'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { PhShieldCheck } from '@phosphor-icons/vue'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({ code: '' })
const rules = {
  code: { required, numeric, minLength: minLength(6), maxLength: maxLength(6) },
}
const v$ = useVuelidate(rules, form)

onMounted(() => {
  if (!authStore.pendingEmail) {
    router.replace('/login')
  }
})

async function handleVerify() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    await authStore.verifyTwoFA(form.value.code)
    swal.success('Login Berhasil')
    router.push('/')
  } catch (error: any) {
    const message = error?.response?.data?.message || 'Kode verifikasi salah atau kedaluwarsa'
    swal.error('Verifikasi Gagal', message)
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <UiCard :padded="false" :classes="{ card: 'p-8' }" class="max-w-md w-full">
      <template #header>
        <h2 class="mb-2 text-center text-2xl font-bold text-gray-800">Verifikasi 2FA</h2>
        <p class="mb-6 text-center text-sm text-gray-500">
          Kode verifikasi telah dikirim ke {{ authStore.pendingEmail }}.
        </p>
      </template>
      <form @submit.prevent="handleVerify">
        <FormInput
          v-model="form.code"
          name="code"
          label="Kode Verifikasi"
          placeholder="123456"
          class="mb-6"
          :validation="v$.code"
          :prefix-icon="PhShieldCheck"
        />

        <UiButton type="submit" full-width :loading="authStore.isLoading" class="mb-4">
          Verifikasi
        </UiButton>

        <div class="text-center text-sm">
          <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
            Kembali ke Login
          </router-link>
        </div>
      </form>
    </UiCard>
  </div>
</template>
