# Environment Variables
<table class="env-vars">
	<tbody>
		<tr>
			<th class="env-vars-section-name" colspan="5">general</th>
		</tr>
		<tr class="env-vars-section-header">
			<th>Переменная</th>
			<th>Тип</th>
			<th>Ограничения</th>
			<th>Значение по умолчанию</th>
			<th>Описание</th>
		</tr>
		<tr>
			<td>CHECK_BATCH_SIZE</td>
			<td><code>int</code></td>
			<td>—</td>
			<td><code>1000</code></td>
			<td>
				<p>Check batch size</p>
			</td>
		</tr>
		<tr>
			<td>CHECK_INTERVAL</td>
			<td><code>time.Duration</code></td>
			<td>—</td>
			<td><code>10m0s</code></td>
			<td>
				<p>Check interval</p>
			</td>
		</tr>
		<tr>
			<td>CHECK_SCHEDULE</td>
			<td><code>string</code></td>
			<td>—</td>
			<td><code>&#34;* * 2-6 * * * *&#34;</code></td>
			<td>
				<p>Check schedule</p>
			</td>
		</tr>
	</tbody>
</table>
