mxpy --verbose contract deploy --bytecode=./output/contract.wasm \
--recall-nonce --pem=wallet-owner.pem \
--gas-limit=100000000 \
--arguments str:BTC-187d25 str:USDT-cf0380 \
--send --outfile="deploy-devnet.interaction.json" --wait-result \
--proxy=https://devnet-gateway.multiversx.com --chain=D


mxpy --verbose contract deploy --bytecode=./output/contract.wasm \
--recall-nonce --pem=wallet-owner.pem \
--gas-limit=100000000 \
--arguments str:ETH-dd05ee str:USDT-cf0380 \
--send --outfile="deploy-devnet.interaction.json" --wait-result \
--proxy=https://devnet-gateway.multiversx.com --chain=D