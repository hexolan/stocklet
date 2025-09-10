<script lang="ts">
	import type { components } from '$lib/api/schema';
	
    const orders: components['schemas']['v1Order'][] = [
        {
            id: '1',
            status: 'ORDER_STATUS_PENDING',
            items: {
                '1': 1,
                '2': 2,
            },
            customerId: '1',
            transactionId: '1',
            createdAt: Date.now().toString(),
            updatedAt: Date.now().toString(),
        },
    ];
</script>